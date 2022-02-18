package lang

var globals = make(map[string]*Object)
var locals = make(map[string]*Object)
var labels = make(map[string][]*Statement)

// DoAll runs as many statements as possible, stopping if there's a problem
// reading the next statement (first value will be false) or if there's a
// problem executing said statement (first value will be true)
func DoAll() error {
	for {
		if !canAdvance() {
			return nil
		}
		stmt, err := NextStmt()
		if err != nil {
			return err
		}
		if stmt == nil {
			continue
		}
		if err = RunStmt(stmt, false); err != nil {
			return err
		}
	}
}

// RunStmt runs a singular statement, consuming more statements if necessary
func RunStmt(stmt *Statement, isLocal bool) error {
	switch stmt.Keyword {
	case "if", "unless":
		if !isLocal { // don't naively wipe locals
			locals = make(map[string]*Object)
		}
		v, err := parseConditional(stmt.Args, stmt.Keyword == "unless")
		if err != nil {
			return err
		}
		for {
			substmt, err := NextStmt()
			if err != nil {
				return err
			}
			if substmt == nil {
				break
			}
			switch substmt.Keyword {
			case "else":
				v = !v
			case "end":
				return nil
			default:
				if v {
					err := RunStmt(substmt, true)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	case "label":
		if isLocal {
			return perr(stmt.KwToken, "labels not allowed in blocks")
		}
		labelName, err := parseObject(stmt.Args)
		if err != nil {
			return err
		}
		if labelName.Type != ObjStr {
			return perr(stmt.Args[0], "label names must be strings")
		}
		labelStmts := []*Statement{}
	labelLoop:
		for {
			substmt, err := NextStmt()
			if err != nil {
				return err
			}
			if substmt == nil {
				break
			}
			switch substmt.Keyword {
			case "end":
				break labelLoop
			default:
				labelStmts = append(labelStmts, substmt)
			}
		}
		labels[labelName.StrV] = labelStmts
		return nil
	case "goto":
		labelName, err := parseObject(stmt.Args)
		if err != nil {
			return err
		}
		if labelName.Type != ObjStr {
			return perr(stmt.Args[0], "label names must be strings")
		}
		labelStmts, ok := labels[labelName.StrV]
		if !ok {
			return perrf(stmt.Args[0], "unknown label %s", labelName.StrV)
		}
		for _, substmt := range labelStmts {
			if err := RunStmt(substmt, true); err != nil {
				return err
			}
		}
		return nil
	case "end":
		return perr(stmt.KwToken, "end statement outside of block")
	case "local", "global", "var", "set":
		name, value, err := parseAssignment(stmt.Args)
		if err != nil {
			return err
		}
		if stmt.Keyword == "local" {
			if !isLocal {
				return perr(stmt.KwToken, "local variable in global context")
			}
			locals[name] = value
		} else if stmt.Keyword == "global" {
			globals[name] = value
		} else {
			if isLocal {
				locals[name] = value
			} else {
				globals[name] = value
			}
		}
		return nil
	default:
		f, ok := Funcs[stmt.Keyword]
		if !ok {
			return perrf(stmt.KwToken, "unknown function %s", stmt.Keyword)
		}
		args, err := parseObjectList(stmt.Args)
		if err != nil {
			return err
		}
		_, err = f(args)
		return err
		// TODO(qeaml): variables, labels and other special statements
	}
}

func GetGlobalVar(name string) (v *Object, ok bool) {
	v, ok = globals[name]
	return
}

func GetLocalVar(name string) (v *Object, ok bool) {
	v, ok = locals[name]
	return
}
