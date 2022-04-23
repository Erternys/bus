package process

func split(command string) []string {
	args := []string{}
	current := ""
	str := false

	for i, c := range command {
		if c == ' ' && !str {
			args = append(args, current)
			current = ""
			continue
		}
		if c == '"' && command[i-1] != '\\' {
			if str {
				args = append(args, current)
				current = ""
			}
			str = !str
			continue
		}
		current += string(c)
	}

	return args
}
