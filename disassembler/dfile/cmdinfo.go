package dfile

import "strings"

type CmdInfo struct {
	MainCmd  string
	Args     []string
	Original string
}

// NewCmdInfos creates CmdInfo per a command.
func NewCmdInfos(s string) []*CmdInfo {
	list := make([]*CmdInfo, 0)
	for _, cmd := range strings.Split(s, "&&") {
		cmd = strings.TrimSpace(cmd)
		cmdInfo := &CmdInfo{
			MainCmd:  "",
			Args:     make([]string, 0),
			Original: cmd,
		}

		for i, arg := range strings.Split(cmd, " ") {
			arg = strings.TrimSpace(arg)
			if arg == "" {
				continue
			}
			if i == 0 {
				cmdInfo.MainCmd = arg
			} else {
				cmdInfo.Args = append(cmdInfo.Args, arg)
			}
		}

		list = append(list, cmdInfo)
	}

	return list
}
