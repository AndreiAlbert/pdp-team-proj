package graph

/*
#cgo CFLAGS: -I/opt/homebrew/include
#include <mpi.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"math"
	"unsafe"
)

func GraphColoringMain(mpiSize int, graph *Graph, colors *Colors) (map[int]string, error) {
	cntColors := colors.CntColors
	codes := graphColoringRec(0, graph, cntColors, make([]int, graph.cntNodes), 0, mpiSize, 0)

	if codes[0] == -1 {
		return nil, errors.New("no solution found")
	}

	return colors.GetNodesToColors(codes), nil
}

func graphColoringRec(node int, graph *Graph, cntColors int, codes []int, mpiId, mpiSize, power int) []int {
	nodesNo := graph.cntNodes

	if !isCodeValid(node, codes, graph) {
		return getArrayOf(nodesNo, -1)
	}

	if node+1 == graph.cntNodes {
		return codes
	}

	coefficient := int(math.Pow(float64(cntColors), float64(power)))
	code := 0
	destination := mpiId + coefficient*(code+1)

	for code < cntColors-1 && destination < mpiSize {
		code++
		destination = mpiId + coefficient*(code+1)
	}

	nextNode := node + 1
	nextPower := power + 1

	for currentCode := 1; currentCode < code; currentCode++ {
		destination = mpiId + coefficient*currentCode

		data := []int{mpiId, nextNode, nextPower}
		C.MPI_Send(unsafe.Pointer(&data[0]), C.int(len(data)), C.MPI_INT, C.int(destination), 0, C.MPI_COMM_WORLD)

		nextCodes := getArrayCopy(codes)
		nextCodes[nextNode] = currentCode

		C.MPI_Send(unsafe.Pointer(&nextCodes[0]), C.int(nodesNo), C.MPI_INT, C.int(destination), 0, C.MPI_COMM_WORLD)

	}

	nextCodes := getArrayCopy(codes)
	nextCodes[nextNode] = 0

	result := graphColoringRec(nextNode, graph, cntColors, nextCodes, mpiId, mpiSize, nextPower)
	if result[0] != -1 {
		return result
	}

	for currentCode := 1; currentCode < code; currentCode++ {
		destination := C.int(mpiId + coefficient*currentCode)

		cResult := (*C.int)(C.malloc(C.size_t(C.sizeof_int * nodesNo)))
		defer C.free(unsafe.Pointer(cResult))

		C.MPI_Recv(unsafe.Pointer(cResult), C.int(nodesNo), C.MPI_INT, destination, C.MPI_ANY_TAG, C.MPI_COMM_WORLD, nil)

		result := make([]int, nodesNo)
		for i := 0; i < nodesNo; i++ {
			result[i] = int(*(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(cResult)) + uintptr(i)*C.sizeof_int)))
		}

		if result[0] != -1 {
			return result
		}
	}

	for currentCode := code; currentCode < cntColors; currentCode++ {
		nextCodes = getArrayCopy(codes)
		nextCodes[nextNode] = currentCode

		result = graphColoringRec(nextNode, graph, cntColors, nextCodes, mpiId, mpiSize, nextPower)
		if result[0] != -1 {
			return result
		}
	}

	return getArrayOf(nodesNo, -1)
}

func GraphColoringWorker(mpiMe, mpiSize int, graph *Graph, codesNo int) {
	nodesNo := graph.cntNodes

	var data [3]C.int
	C.MPI_Recv(unsafe.Pointer(&data), C.int(len(data)), C.MPI_INT, C.MPI_ANY_SOURCE, C.MPI_ANY_TAG, C.MPI_COMM_WORLD, nil)

	parent := data[0]
	node := data[1]
	power := data[2]

	cCodes := (*C.int)(C.malloc(C.size_t(C.sizeof_int) * C.size_t(nodesNo)))
	defer C.free(unsafe.Pointer(cCodes))

	C.MPI_Recv(unsafe.Pointer(cCodes), C.int(nodesNo), C.MPI_INT, C.MPI_ANY_SOURCE, C.MPI_ANY_TAG, C.MPI_COMM_WORLD, nil)

	goCodes := make([]int, nodesNo)
	for i := 0; i < nodesNo; i++ {
		goCodes[i] = int(*(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(cCodes)) + uintptr(i)*uintptr(C.sizeof_int))))
	}

	newCodes := graphColoringRec(int(node), graph, codesNo, goCodes, mpiMe, mpiSize, int(power))

	cNewCodes := (*C.int)(C.malloc(C.size_t(C.sizeof_int) * C.size_t(nodesNo)))
	defer C.free(unsafe.Pointer(cNewCodes))

	for i, v := range newCodes {
		*(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(cNewCodes)) + uintptr(i)*uintptr(C.sizeof_int))) = C.int(v)
	}

	C.MPI_Send(unsafe.Pointer(cNewCodes), C.int(nodesNo), C.MPI_INT, C.int(parent), 0, C.MPI_COMM_WORLD)
}

func isCodeValid(node int, codes []int, graph *Graph) bool {
	for currentNode := 0; currentNode < node; currentNode++ {
		if (graph.IsEdge(node, currentNode) || graph.IsEdge(currentNode, node)) && codes[node] == codes[currentNode] {
			return false
		}
	}
	return true
}

func getArrayOf(length, value int) []int {
	array := make([]int, length)
	for i := range array {
		array[i] = value
	}
	return array
}

func getArrayCopy(array []int) []int {
	newArray := make([]int, len(array))
	copy(newArray, array)
	return newArray
}
