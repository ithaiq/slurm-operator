package controllers

import "fmt"

type (
	SlurmResourceDefineName = string
	SlurmResourceType       = string
)

const (
	SlurmApplicationLabelKey                         = "slurmoperator.apulis.cn/cluster"
	SlurmJupyter             SlurmResourceDefineName = "jupyter"
	SlurmMaster              SlurmResourceDefineName = "master"
	SlurmNode                SlurmResourceDefineName = "node"

	SlurmPod     SlurmResourceType = "pod"
	SlurmService SlurmResourceType = "service"
)

type Command = string

const (
	SedMasterName Command = "sed -i 's/slurmmaster/%s/g' /etc/slurm-llnl/slurm.conf"
	SedNodeName   Command = "sed -i 's/slurmnode/%s/g' /etc/slurm-llnl/slurm.conf"
	SedMLogLevel  Command = "sed -i 's/SlurmctldDebug=error/SlurmctldDebug=verbose/g' /etc/slurm-llnl/slurm.conf"
	SedNLogLevel  Command = "sed -i 's/SlurmdDebug=error/SlurmdDebug=verbose/g' /etc/slurm-llnl/slurm.conf"
	SedSLogLevel  Command = "sed -i 's/#SlurmSchedLogLevel=/SlurmSchedLogLevel=1/g' /etc/slurm-llnl/slurm.conf"
	SedMLogFile   Command = "sed -i 's/SlurmctldLogFile=\\/var\\/log\\/slurm-llnl\\/slurmctld.log/SlurmctldLogFile=\\/dev\\/stdout/g' /etc/slurm-llnl/slurm.conf"
	SedNLogFile   Command = "sed -i 's/SlurmdLogFile=\\/var\\/log\\/slurm-llnl\\/slurmd.log/SlurmdLogFile=\\/dev\\/stdout/g' /etc/slurm-llnl/slurm.conf"
	SedSLogFile   Command = "sed -i 's/#SlurmSchedLogFile=/SlurmSchedLogFile=\\/dev\\/stdout/g' /etc/slurm-llnl/slurm.conf"
	CmdRun        Command = "/etc/slurm-llnl/docker-entrypoint.sh"
)

func GetSlurmResourceName(slurmName string, defineName SlurmResourceDefineName) string {
	return slurmName + "-" + defineName
}

func CombineSlurmResourceType(resourceName string, resourceType SlurmResourceType) string {
	return resourceName + "-" + resourceType
}

func GetPrePodRunCommand(cmdList ...Command) string {
	var cmd string
	for index, cmdItem := range cmdList {
		if index == len(cmdList)-1 {
			cmd += fmt.Sprintf(cmdItem)
			break
		}
		cmd += fmt.Sprintf(cmdItem + " && ")
	}
	return cmd
}
