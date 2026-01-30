package generator

import "fmt"

// Label names for control flow.
func whileStartLabel(id int) string    { return fmt.Sprintf("while_start_%d", id) }
func whileEndLabel(id int) string      { return fmt.Sprintf("while_end_%d", id) }
func loopContinueLabel(id int) string  { return fmt.Sprintf("loop_continue_%d", id) }
func loopBreakLabel(id int) string     { return fmt.Sprintf("loop_break_%d", id) }
func fnReturnLabel(name string) string { return fmt.Sprintf("fn_ret_%s", name) }

// Mangle function names deterministically.
func mangleFunc(name string) string { return fmt.Sprintf("fn_%s", name) }

// Mangle temporary variable names deterministically.
func mangleTemp(prefix string, id int) string { return fmt.Sprintf("%s_tmp_%d", prefix, id) }
