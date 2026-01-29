package generator

import "fmt"

// Label names for control flow.
func whileStartLabel(id int) string { return fmt.Sprintf("while_start_%d", id) }
func whileEndLabel(id int) string   { return fmt.Sprintf("while_end_%d", id) }

// Mangle function names deterministically.
func mangleFunc(name string) string { return fmt.Sprintf("fn_%s", name) }

// Mangle temporary variable names deterministically.
func mangleTemp(prefix string, id int) string { return fmt.Sprintf("%s_tmp_%d", prefix, id) }
