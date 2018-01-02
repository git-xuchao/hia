package contract

import (
	"fmt"
	"testing"
)

func Test_DeployContractPipeline(t *testing.T) {
	fmt.Println("begin to deploy contract\n")
	DeployContractPipeline("12310", "/home/seagull/projects/contracts/copyright/copyright.sol")
}

func Test_ExecCompileContract(t *testing.T) {
	fmt.Println("begin to compile contract\n")
	ExecCompileContract("/home/seagull/projects/contracts/copyright/copyright.sol")
}

func Test_DeployContract(t *testing.T) {
	fmt.Println("begin to compile contract\n")
	DeployContract("12312", "_home_seagull_projects_contracts_copyright_copyright_sol_copyright.abi", "_home_seagull_projects_contracts_copyright_copyright_sol_copyright.bin")
}
