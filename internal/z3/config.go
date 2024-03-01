package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"
import "runtime"

// Configuration object used to initialize logical contexts.
type Config struct {
	z3Config C.Z3_config
}

// Create a configuration object for the Z3 context object.
// Configurations are created in order to assign parameters
// prior to creating contexts for Z3 interaction.
func NewConfig() *Config {
	config := &Config{
		// The following parameters can be set:
		//     - proof  (Boolean)           Enable proof generation
		//     - debug_ref_count (Boolean)  Enable debug support for Z3_ast reference counting
		//     - trace  (Boolean)           Tracing support for VCC
		//     - trace_file_name (String)   Trace out file for VCC traces
		//     - timeout (unsigned)         default timeout (in milliseconds) used for solvers
		//     - well_sorted_check          type checker
		//     - auto_config                use heuristics to automatically select solver and configure it
		//     - model                      model generation for solvers, this parameter can be overwritten when creating a solver
		//     - model_validate             validate models produced by solvers
		//     - unsat_core                 unsat-core generation for solvers, this parameter can be overwritten when creating a solver
		//     - encoding                   the string encoding used internally (must be either "unicode" - 18 bit, "bmp" - 16 bit or "ascii" - 8 bit)
		z3Config: C.Z3_mk_config(),
	}

	// We have to delete the config after the construction of the Context.
	//   Otherwise, we would have a memory leak.
	runtime.SetFinalizer(config, func(config *Config) {
		C.Z3_del_config(config.z3Config)
	})

	return config
}
