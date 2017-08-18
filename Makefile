COMMIT = $$(git describe --always)     

PKG_FILES ?= *.go                      
LINT_LOG := lint.log                   

# Run lint on the latest stable release.                                       
GO_VERSION := $(shell go version | cut -d " " -f 3)                            
GO_MINOR_VERSION := $(word 2,$(subst ., ,$(GO_VERSION)))                       
LINTABLE_MINOR_VERSIONS := 8           
ifneq ($(filter $(LINTABLE_MINOR_VERSIONS),$(GO_MINOR_VERSION)),)              
	SHOULD_LINT := true            
endif                                  

BENCH_FLAGS ?= -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem           

deps:                                  
	@echo "Installing Glide and dependencies..."                           
	glide --version || go get -u -f github.com/Masterminds/glide           
	glide install                  

lint:                                  
ifdef SHOULD_LINT                      
	@rm -rf $(LINT_LOG)            
	@echo "Checking format..."     
	@gofmt -d -s $(PKG_FILES) 2>&1 | tee $(LINT_LOG)                       
	@echo "Checking vet..."        
	@$(foreach dir,$(PKG_FILES),go tool vet $(dir) 2>&1 | tee -a $(LINT_LOG))                                                                             
	@echo "Running golint..."      
	@golint ./... 2>&1 | tee -a $(LINT_LOG)                                
else                                   
	@echo "Skipping linters"       
endif                                  

test:                                  
	@echo "Running test..."        
	go test -race                  

bench:                                 
	@echo "Running benchmark test..."                                      
	go test -run="^$$" $(BENCH_FLAGS)                                      

.PHONY: deps lint test bench           
