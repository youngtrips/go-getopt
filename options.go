package getopt

type Options []Option;

func (options Options) FindOption(optionString string) (option Option, found bool) {
  for _, cur := range options {
    if cur.ShortOpt() == optionString || cur.LongOpt() == optionString {
      option = cur
      found = true
      break
    }
  }

  return option, found
}

func (options Options) IsOptional(optionName string) (isRequired bool) {
  if option, found := options.FindOption(optionName); found && option.flags & Optional != 0 {
    isRequired = true
  }

  return isRequired
}

func (options Options) IsRequired(optionName string) (isRequired bool) {
  if option, found := options.FindOption(optionName); found && option.flags & Required != 0 {
    isRequired = true
  }

  return isRequired
}

func (options Options) IsFlag(optionName string) (isFlag bool) {
  if option, found := options.FindOption(optionName); found && option.flags & Flag != 0 {
    isFlag = true
  }

  return isFlag
}

func (options Options) RequiredOptions() (requiredOptions []string) {

  for _, cur := range options {
    if cur.flags & Required != 0 {
      requiredOptions = append(requiredOptions, cur.LongOpt())
    }
  }

  return requiredOptions
}

func (options Options) Usage(programName string) (output string) {
  output = "\n\n    Usage: " + programName

  for _, option := range options {
    output = output + " " + option.Usage()
  }

  output = output + "\n\n"

  return
}

func (options Options) Help(programName string, description string) (output string) {
  output = options.Usage(programName)
  if description != "" {
    output = output + description + "\n"
  }

  longOptTextLength := 0

  for _, option := range options {
    if length := len(option.LongOptString()); length > longOptTextLength {
      longOptTextLength = length
    }
  }

  longOptTextLength = longOptTextLength + 2

  var argumentsString string
  var optionsString string

  for _, option := range options {
    if option.flags & IsArg > 0 {
      argumentsString = argumentsString + option.HelpText(longOptTextLength) + "\n"
    } else {
      optionsString = optionsString + option.HelpText(longOptTextLength) + "\n"
    }
  }

  if optionsString != "" {
    output = output + "\nOptions:\n" + optionsString
  }

  if argumentsString != "" {
    output = output + "\nArguments:\n" + argumentsString + "\n"
  }

  return
}
