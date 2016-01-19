package metrix

import (
	"strconv"
	"time"
)

type VmState int
type LcmState int

const (
	//VmState starts at 0
	INIT VmState = iota
	PENDING
	HOLD
	ACTIVE
	STOPPED
	SUSPENDED
	DONE
	FAILED

	//LcmState starts at 0
	LCM_INIT LcmState = iota
	PROLOG
	BOOT
	RUNNING
	MIGRATE
	SAVE_STOP
	SAVE_SUSPEND
	SAVE_MIGRATE
	PROLOG_MIGRATE
	PROLOG_RESUME
	EPILOG_STOP
	EPILOG
	SHUTDOWN
	CANCEL
	FAILURE
	CLEANUP
	UNKNOWN
)

func timeAsInt64(tm string) int64 {
	if i, err := strconv.ParseInt(tm, 10, 64); err != nil {
		return i
	}
	return 0
}

type OneContext struct {
	NAME          string `json:"NAME"`
	ACCOUNTS_ID   string `json:"ACCOUNTS_ID"`
	ASSEMBLY_ID   string `json:"ASSEMBLY_ID"`
	ASSEMBLIES_ID string `json:"ASSEMBLIES_ID"`
}

type OneTemplate struct {
	CONTEXT     OneContext `json:"CONTEXT"`
	CPU         string     `json:"CPU"`
	CPU_COST    string     `json:"CPU_COST"`
	VCPU        string     `json:"VCPU"`
	MEMORY      string     `json:"MEMORY"`
	MEMORY_COST string     `json:"MEMORY_COST"`
	DISK_SIZE   string     `json:"SIZE"`
}

type OneVM struct {
	NAME      string       `json:"NAME"`
	STATE     string       `json:"STATE"`
	LCM_STATE string       `json:"LCM_STATE"`
	STIME     string       `json:"STIME"`
	ETIME     string       `json:"ETIME"`
	TEMPLATE  *OneTemplate `json:"TEMPLATE"`
}

//Returns the delta duration in hours as string
func (ov *OneVM) elapsed() string {
	elapsed := strconv.FormatFloat(time.Since(time.Unix(timeAsInt64(ov.STIME), 0)).Hours(), 'E', -1, 64)
	return elapsed
}

//Returns the delta duration in hours as string
func (ov *OneVM) stateAsInt(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return 22
}

type OneHistory struct {
	HOSTNAME string `json:"HOSTNAME"`
	STIME    string `json:"STIME"`
	ETIME    string `json:"ETIME"`
	VM       *OneVM `json:"VM"`
}

type OpenNebulaStatus struct {
	HISTORYS []*OneHistory `json:"HISTORY"`
}

func (oh *OneHistory) Cpu() string {
	return oh.VM.TEMPLATE.CPU
}

func (oh *OneHistory) CpuCost() string {
	return oh.VM.TEMPLATE.CPU_COST
}

func (oh *OneHistory) Memory() string {
	return oh.VM.TEMPLATE.MEMORY
}

func (oh *OneHistory) MemoryCost() string {
	return oh.VM.TEMPLATE.MEMORY_COST
}

func (oh *OneHistory) AssemblyName() string {
	return oh.VM.NAME
}

func (oh *OneHistory) AccountsId() string {
	return oh.VM.TEMPLATE.CONTEXT.ACCOUNTS_ID
}

func (oh *OneHistory) AssembliesId() string {
	return oh.VM.TEMPLATE.CONTEXT.ASSEMBLIES_ID

}

func (oh *OneHistory) AssemblyId() string {
	return oh.VM.TEMPLATE.CONTEXT.ASSEMBLY_ID
}

func (oh *OneHistory) State() string {
	return oh.VM.stateString()
}

func (oh *OneHistory) LcmState() string {
	return oh.VM.lcmStateString()
}

func (ov *OneVM) stateString() string {
	switch VmState(ov.stateAsInt(ov.STATE)) {
	case INIT:
		return "Init"
	case PENDING:
		return "Pending"
	case HOLD:
		return "Hold"
	case ACTIVE:
		return "Active"
	case STOPPED:
		return "Stopped"
	case SUSPENDED:
		return "Suspended"
	case DONE:
		return "Done"
	case FAILED:
		return "Failed"
	default:
	}
	return "Unknown"
}

func (ov *OneVM) lcmStateString() string {
	switch LcmState(ov.stateAsInt(ov.LCM_STATE)) {
	case LCM_INIT:
		return "Lcm Init"
	case PROLOG:
		return "Prolog"
	case BOOT:
		return "Boot"
	case RUNNING:
		return "Running"
	case MIGRATE:
		return "Migrate"
	case SAVE_STOP:
		return "Save stop"
	case SAVE_SUSPEND:
		return "Save suspend"
	case SAVE_MIGRATE:
		return "Save migrate"
	case PROLOG_MIGRATE:
		return "Prolog migrate"
	case PROLOG_RESUME:
		return "Prolog resume"
	case EPILOG_STOP:
		return "Eplilog stop"
	case EPILOG:
		return "Epilog"
	case SHUTDOWN:
		return "Shutdown"
	case CANCEL:
		return "Cancel"
	case FAILURE:
		return "Failure"
	case CLEANUP:
		return "Cleanup"
	default:
		return "Unknown"
	}
}
