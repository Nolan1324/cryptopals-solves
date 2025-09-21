package timeattack

import (
	"context"
	"log"
	"os/exec"
	"sync"
)

// RunServerProcess starts the server in a seperate process and then blocks
func RunServerProcess(ctx context.Context, mainGoPath string) {
	cmd := exec.CommandContext(ctx, "go", "run", mainGoPath)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Println("waiting on cmd: ", err)
	}
}

// RunServerProcess starts the server in a seperate process without blocking,
// and returns a function to stop the server
func RunServerProcessInBackground(mainGoPath string) func() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		RunServerProcess(ctx, mainGoPath)
	}()

	return func() {
		cancel()
		wg.Wait()
	}
}
