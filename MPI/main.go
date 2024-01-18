package MPI

/*
   #cgo CFLAGS: -I/opt/homebrew/include
   #cgo LDFLAGS: -L/opt/homebrew/lib -lmpi
   #include <mpi.h>
*/
import "C"
import (
	"fmt"
	"github.com/AndreiAlbert/n-coloring-graph/MPI/graph"
	"log"
	"strings"
)

func PrettyPrint(coloring map[int]string) string {
	if coloring == nil || len(coloring) == 0 {
		return "No coloring available."
	}

	var builder strings.Builder
	builder.WriteString("Graph Coloring Result:\n")

	for node, color := range coloring {
		builder.WriteString(fmt.Sprintf("Node %d: %s\n", node, color))
	}

	return builder.String()
}

func main() {
	C.MPI_Init(nil, nil)
	defer C.MPI_Finalize()

	var rank, size C.int
	C.MPI_Comm_rank(C.MPI_COMM_WORLD, &rank)
	C.MPI_Comm_size(C.MPI_COMM_WORLD, &size)

	nodesNo := 5
	g := graph.NewGraph(nodesNo)
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 4)
	g.AddEdge(4, 0)
	g.AddEdge(2, 0)
	g.AddEdge(0, 4)
	g.AddEdge(4, 3)
	g.AddEdge(3, 1)

	colors := graph.NewColors(3)
	colors.SetColorName(0, "red")
	colors.SetColorName(1, "green")
	colors.SetColorName(2, "blue")

	if rank == 0 {
		fmt.Println("Master process")

		result, err := graph.GraphColoringMain(int(size), g, colors)

		if err != nil {
			log.Fatalf("Error: %s", err.Error())
			return
		}

		fmt.Println(PrettyPrint(result))
	} else {
		fmt.Println("Worker process:", rank)

		codesNo := colors.CntColors

		graph.GraphColoringWorker(int(rank), int(size), g, codesNo)
	}
}
