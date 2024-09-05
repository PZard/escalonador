package main

type Process struct {
	name     string
	cpuBurst int
	ioTime   int
	totalCpu int
	priority int
	credits  int
	state    string
	order    int
}

var processList []Process

func scheduler(processes []Process) {
	for {
		readyProcesses := getReadyProcesses(processes)
		if len(readyProcesses) == 0 {
			reassignCredits(processes)
			continue
		}

		process := selectHighestCreditProcess(readyProcesses)

		simulateCpuExecution(&process)

		if allProcessesFinished(processes) {
			break
		}
	}
}

func getReadyProcesses(processes []Process) []Process {
	var ready []Process
	for _, p := range processes {
		if p.state == "Ready" {
			ready = append(ready, p)
		}
	}
	return ready
}

func selectHighestCreditProcess(processes []Process) Process {
	highestCredit := processes[0]
	for _, p := range processes {
		if p.credits > highestCredit.credits ||
			(p.credits == highestCredit.credits && p.priority > highestCredit.priority) {
			highestCredit = p
		}
	}
	return highestCredit
}

func simulateCpuExecution(process *Process) {
	if process.cpuBurst > 0 {
		process.cpuBurst--
		process.credits--
		if process.cpuBurst == 0 {
			process.state = "Blocked"
			process.ioTime = 5 // Example of I/O time
		}
	} else if process.ioTime > 0 {
		process.ioTime--
		if process.ioTime == 0 {
			process.state = "Ready"
		}
	}
}

func reassignCredits(processes []Process) {
	for i := range processes {
		processes[i].credits = (processes[i].credits / 2) + processes[i].priority
		if processes[i].state == "Blocked" {
			processes[i].state = "Ready"
		}
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
	processList = []Process{
		{"A", 2, 5, 6, 3, 3, "Ready", 1},
		{"B", 3, 10, 6, 3, 3, "Ready", 2},
		{"C", 14, 0, 14, 3, 3, "Ready", 3},
		{"D", 10, 0, 10, 3, 3, "Ready", 4},
	}

	scheduler(processList)
}
