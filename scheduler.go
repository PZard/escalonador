package main

type Process struct {
	name             string
	originalCpuBurst int
	cpuBurst         int
	originalIoTime   int
	ioTime           int
	totalCpu         int
	order            int
	priority         int
	credits          int
	state            string
}

func NewProcess(name string, cpuBurst int, ioTime int, totalCpu int, order int, priority int) Process {
	return Process{
		name:             name,
		originalCpuBurst: cpuBurst,
		cpuBurst:         cpuBurst,
		ioTime:           ioTime,
		originalIoTime:   ioTime,
		totalCpu:         totalCpu,
		order:            order,
		priority:         priority,
		credits:          priority,
		state:            "Ready",
	}
}

func scheduler(processes []Process) {
	var currentProcess *Process
	var oldProcess *Process

	for {
		if currentProcess == nil || currentProcess.state == "Blocked" || currentProcess.credits == 0 {
			readyProcesses := getReadyProcesses(processes)

			if len(readyProcesses) == 0 {
				reassignCredits(processes)
				readyProcesses = getReadyProcesses(processes)
			}

			currentProcess = selectHighestCreditProcess(readyProcesses, oldProcess)
		}

		blockedCurrent := simulateCpuExecution(currentProcess)

		waitBlockeds(processes, currentProcess, blockedCurrent)

		if blockedCurrent {
			currentProcess.order = len(processes) + currentProcess.order
		}

		if currentProcess.credits == 0 {
			oldProcess = currentProcess
			oldProcess.order = len(processes) + oldProcess.order
			currentProcess = nil
			continue
		}

		if currentProcess.state == "Exit" {
			oldProcess = currentProcess
			oldProcess.order = len(processes) + oldProcess.order
			currentProcess = nil
		}

		if allProcessesFinished(processes) {
			break
		}
	}
}

func getReadyProcesses(processes []Process) []*Process {
	var ready []*Process
	for i := range processes {
		if processes[i].state == "Ready" && processes[i].credits > 0 {
			ready = append(ready, &processes[i])
		}
	}
	return ready
}

func selectHighestCreditProcess(processes []*Process, oldProcess *Process) *Process {
	var highestCredit *Process
	for _, p := range processes {
		if oldProcess != nil && p.name == oldProcess.name {
			continue
		}
		if highestCredit == nil {
			highestCredit = p
			continue
		}
		if p.credits > highestCredit.credits || (p.credits == highestCredit.credits && p.order < highestCredit.order) {
			highestCredit = p
		}
	}
	return highestCredit
}

func simulateCpuExecution(process *Process) bool {
	blocked := false
	if process.cpuBurst > 0 {
		process.cpuBurst--

		if process.cpuBurst == 0 {
			process.cpuBurst = process.originalCpuBurst
			process.state = "Blocked"
			blocked = true
		}
	}

	process.totalCpu--
	process.credits--

	if process.totalCpu == 0 {
		process.state = "Exit"
	}

	return blocked;
}

func waitBlockeds(processes []Process, current *Process, blocked bool) {
	for i := range processes {
		if blocked && current.name == processes[i].name {
			continue
		}
		if processes[i].state == "Blocked" && processes[i].ioTime > 0 {
			processes[i].ioTime--

			if processes[i].ioTime == 0 {
				processes[i].ioTime = processes[i].originalIoTime
				processes[i].state = "Ready"
			}
		}
	}
}

func reassignCredits(processes []Process) {
	for i := range processes {
		processes[i].credits = (processes[i].credits / 2) + processes[i].priority
	}
}

func allProcessesFinished(processes []Process) bool {
	for _, p := range processes {
		if p.state != "Exit" {
			return false
		}
	}
	return true
}

func main() {
	processList := []Process{
		NewProcess("A", 2, 5, 6, 1, 3),
		NewProcess("B", 3, 10, 6, 2, 3),
		NewProcess("C", 0, 0, 14, 3, 3),
		NewProcess("D", 0, 0, 10, 4, 3),
	}

	scheduler(processList)
}
