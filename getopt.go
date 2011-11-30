package getopt

const Required = 1
const Optional = 2
const Flag = 3

const InvalidOption = 100
const NoShortOpt    = 101

const OPTIONS_SEPARATOR = "--"

type option struct {
  long_opt string
  short_opt string
  arg_type int
  description string
  default_value interface{}
}

type Options []option;
type GetOptError struct {
  errorCode int
  message string
}

func isShortOpt(option string) (opt string, val string, found bool) {
  if len(option) > 1 && option[0] == '-' && option[1] >= 'A' && option[1] <= 'z' {
    found = true
    opt = option[1:2]
    if len(option) > 2 {
      val = option[2:]
    }

  }

  return opt, val, found
}

func isLongOpt(option string) {

}

func parseShortOpt(shortOpt string, optionDefinition Options) (key string, value string, err *GetOptError) {
  if key[0] != '-' { return "", "", &GetOptError{NoShortOpt, ""} }

  return "", "", nil
}

func (optionsDefinition Options) parse(args []string) (map[string] string, []string, []string, *GetOptError) {
  options := make(map[string] string)
  requiredValues := make([]string, 0)

  for _, option := range optionsDefinition {
    if option.arg_type == Required {
      requiredValues = append(requiredValues, option.long_opt)
    }
  }

  for i:=1; i<len(args); i++ { // args[0] is no opt
  }

  for _, option := range requiredValues {
    if _, found := options[option]; found == false {
      return options, []string{}, []string{}, &GetOptError{1, "required option " + option + " is not set"}
    }
  }

  return options, []string{}, []string{}, nil
}
