APP=bin/app
PPROF_HOST=localhost
PPROF_PORT=8081
PROFILE_DURATION=30

prof-cpu: ## Run CPU profile
	go tool pprof -output=cpu.prof "http://$(PPROF_HOST):$(PPROF_PORT)/debug/profile?seconds=$(PROFILE_DURATION)"

prof-heap: ## Run memory profile
	go tool pprof -output=heap.prof "http://$(PPROF_HOST):$(PPROF_PORT)/debug/heap"

prof-goroutine: ## Run goroutine profile
	go tool pprof -output=goroutine.prof "http://$(PPROF_HOST):$(PPROF_PORT)/debug/goroutine"

prof-block: ## Run block profile
	go tool pprof -output=block.prof "http://$(PPROF_HOST):$(PPROF_PORT)/debug/block"

prof-mutex: ## Run locks profile
	go tool pprof -output=mutex.prof "http://$(PPROF_HOST):$(PPROF_PORT)/debug/mutex"

prof-viz-cpu: ## Visualize CPU profile
	go tool pprof -http=:8080 $(APP) cpu.prof

prof-viz-heap: ## Visualize memory profile
	go tool pprof -http=:8080 $(APP) heap.prof

prof-viz-goroutine: ## Visualize goroutine profile
	go tool pprof -http=:8080 $(APP) goroutine.prof

prof-viz-block: ## Visualize block profile
	go tool pprof -http=:8080 $(APP) block.prof

prof-viz-mutex: ## Visualize locks profile
	go tool pprof -http=:8080 $(APP) mutex.prof

prof-clean: ## Clean *.prof files
	rm -f $(APP) *.prof
