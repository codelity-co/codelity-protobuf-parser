package parser

import (
	"errors"
	"fmt"
	"strings"
)

type Parser struct {
	syntax      string
	packageName string
	options     map[string]string
	imports     []string
	services    []*Service
	messages    []*Message
}

// func init() {
// 	var err error
// 	buildConfig := zap.NewProductionConfig()
// 	logger, err = buildConfig.Build()
// 	if err != nil {
// 		panic(err)
// 	}
// 	sugar = logger.Sugar()
// }

// Parse will accept content as string and translate content into type Parser
func (p *Parser) Parse(content string) error {
	lines := strings.Split(content, "\n")

	syntaxLines := make([]string, 0)
	packageLines := make([]string, 0)
	optionLines := make([]string, 0)
	importLines := make([]string, 0)
	serviceLines := make([]string, 0)
	messageLines := make([]string, 0)

	lastOperation := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		tokens := strings.Split(line, " ")
		switch tokens[0] {
		case "syntax":
			lastOperation = "syntax"
			syntaxLines = append(syntaxLines, line)
		case "package":
			lastOperation = "package"
			packageLines = append(packageLines, line)
		case "rpc":
			lastOperation = "rpc"
			serviceLines = append(serviceLines, line)
		case "option":
			if lastOperation == "rpc" {
				serviceLines = append(serviceLines, line)
			} else {
				lastOperation = "option"
				optionLines = append(optionLines, line)
			}
		case "import":
			lastOperation = "import"
			importLines = append(importLines, line)
		case "service":
			lastOperation = "service"
			serviceLines = append(serviceLines, line)
		case "message":
			lastOperation = "message"
			messageLines = append(messageLines, line)
		default:
			switch lastOperation {
			case "syntax":
				syntaxLines = append(syntaxLines, line)
			case "package":
				packageLines = append(packageLines, line)
			case "rpc":
				serviceLines = append(serviceLines, line)
			case "option":
				if lastOperation == "rpc" {
					serviceLines = append(serviceLines, line)
				} else {
					optionLines = append(optionLines, line)
				}
			case "import":
				importLines = append(importLines, line)
			case "service":
				serviceLines = append(serviceLines, line)
			case "message":
				messageLines = append(messageLines, line)
			}
		}

	}

	err := p.processSyntaxLines(syntaxLines)
	if err != nil {
		return err
	}

	err = p.processPackageLines(packageLines)
	if err != nil {
		return err
	}

	err = p.processImportLines(importLines)
	if err != nil {
		return err
	}

	err = p.processOptionLines(optionLines)
	if err != nil {
		return err
	}

	err = p.processServiceLines(serviceLines)
	if err != nil {
		return err
	}

	err = p.processMessageLines(messageLines)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) ReadFile(filePath string) error {
	return nil
}

func (p *Parser) GetSyntax() string {
	return p.syntax
}

func (p *Parser) GetPackageName() string {
	return p.packageName
}

func (p *Parser) GetOptions() map[string]string {
	return p.options
}

func (p *Parser) GetImports() []string {
	return p.imports
}

func (p *Parser) GetServices() []*Service {
	return p.services
}

func (p *Parser) GetMessages() []*Message {
	return p.messages
}

func (p *Parser) processSyntaxLines(lines []string) error {
	//nolint: errcheck
	// defer sugar.Sync()
	var err error

	// sugar.Debugf("syntaxLines: %v", lines)

	if len(lines) == 0 {
		err = errors.New(`Missing syntax declaration`)
		// sugar.Error(err)
		return err
	}

	if len(lines) > 1 {
		err = errors.New(`Multiple syntax lines`)
		// sugar.Error(err)
		return err
	}

	tokens := strings.Split(lines[0], "=")
	if len(tokens) < 2 {
		err = fmt.Errorf(`syntax line does not have equal sign: %v`, lines[0])
		// sugar.Error(err)
		return err
	}

	semiColonPos := strings.Index(tokens[1], `;`)
	if semiColonPos == -1 {
		err = fmt.Errorf(`syntax line is not ending with semicolon: %v`, lines[0])
		// sugar.Error(err)
		return err
	}

	syntax := strings.Split(tokens[1], ";")
	value := strings.TrimSpace(strings.Replace(strings.Replace(syntax[0], `"`, ``, -1), `'`, ``, -1))
	if value != "proto3" {
		err = fmt.Errorf(`syntax is not "proto3": %v`, lines[0])
		// sugar.Error(err)
		return err
	}

	p.syntax = value
	return nil
}

func (p *Parser) processPackageLines(lines []string) error {
	//nolint: errcheck
	// defer sugar.Sync()
	var err error

	// sugar.Debugf("packageLines: %v", lines)

	if len(lines) == 0 {
		err = errors.New(`Missing package declaration`)
		// sugar.Error(err)
		return err
	}

	if len(lines) > 1 {
		err = errors.New("Multiple package lines")
		// sugar.Error(err)
		return err
	}

	line := strings.Join(lines, " ")
	semiColonPos := strings.Index(line, `;`)
	if semiColonPos == -1 {
		err = fmt.Errorf(`Package line is not ending with semicolon: %v`, line)
		// sugar.Error(err)
		return err
	}

	dobuleQuotePos := strings.Index(line, `"`)
	if dobuleQuotePos > -1 {
		err = fmt.Errorf(`Package line should not have double quote: %v`, line)
		// sugar.Error(err)
		return err
	}

	singleQuotePos := strings.Index(line, `'`)
	if singleQuotePos > -1 {
		err = fmt.Errorf(`Package line should not have single quote: %v`, line)
		// sugar.Error(err)
		return err
	}

	packageName := strings.TrimSpace(strings.Replace(strings.Replace(line, `package`, ``, -1), `;`, ``, -1))
	if len(packageName) == 0 {
		err = fmt.Errorf(`Missing package name: %v`, line)
		return err
	}

	p.packageName = packageName
	return nil
}

func (p *Parser) processOptionLines(lines []string) error {
	//nolint: errcheck
	// defer sugar.Sync()

	var err error
	// sugar.Debugf("optionLines: %v", lines)

	p.options = make(map[string]string)

	if len(lines) == 0 {
		return nil
	}

	var line string
	if len(lines) > 1 {
		line = strings.Join(lines, "\n")
	} else if len(lines) == 1 {
		line = lines[0] + "\n"
	}

	numberOfOptions := strings.Count(line, "option")
	numberOfSemicolon := strings.Count(line, ";\n")

	if numberOfOptions > numberOfSemicolon {
		err = fmt.Errorf("Missing semicolon is found in option statements: %v", line)
		return err
	}

	if numberOfOptions < numberOfSemicolon {
		err = fmt.Errorf("Dummy semicolon is found in option statements: %v", line)
		return err
	}

	lines = strings.Split(line, ";\n")

	for _, line := range lines {

		tempLine := strings.TrimSpace(line)
		if len(tempLine) == 0 {
			continue
		}

		doubleQuoteCount := strings.Count(tempLine, `"`)
		singleQuoteCount := strings.Count(tempLine, `'`)

		if doubleQuoteCount == 0 && singleQuoteCount == 0 {
			err = fmt.Errorf(`Missing quote in option statements: %v`, tempLine)
			return err
		}

		if (doubleQuoteCount % 2) == 1 || (singleQuoteCount % 2) == 1 {
			fmt.Println(tempLine)
			err = fmt.Errorf(`Mismatch quotes are found in option statements: %v`, tempLine)
			return err
		}

		if doubleQuoteCount > 0 && strings.LastIndex(line, `"`) != (len(tempLine) - 1) {
			err = fmt.Errorf(`Invalid option statements: %v`, tempLine)
			return err
		}

		if singleQuoteCount > 0 && strings.LastIndex(line, `'`) != (len(tempLine) - 1) {
			err = fmt.Errorf(`Invalid option statements: %v`, tempLine)
			return err
		}

		tokens := strings.Split(tempLine, ` `)
		if len(tokens) != 3 {
			err = fmt.Errorf(`Invalid option statements: %v`, tempLine)
			return err
		}

		optionKey := tokens[1]
		if strings.Count(optionKey, `"`) > 0 || strings.Count(optionKey, `'`) > 0 {
			err = fmt.Errorf(`Invalid option key: %v`, tempLine)
			return err
		}

		optionValue := strings.ReplaceAll(strings.ReplaceAll(tokens[2], `"`, ``), `'`, ``)
		if len(optionValue) == 0 {
			err = fmt.Errorf(`Invalid option value: %v`, tempLine)
			return err
		}

		p.options[optionKey] = optionValue
	}

	return nil
}

func (p *Parser) processImportLines(lines []string) error {
	//nolint: errcheck
	// defer sugar.Sync()

	var err error
	// sugar.Debugf("importLines: %v", lines)

	p.imports = make([]string, 0)

	if len(lines) == 0 {
		return nil
	}

	line := strings.Join(lines, " ")

	numberOfImports := strings.Count(line, "import")
	numberOfSemicolon := strings.Count(line, ";")

	if numberOfImports > numberOfSemicolon {
		err = fmt.Errorf(`Missing semicolon is found in import statements: %v`, line)
		return err
	}

	if numberOfImports < numberOfSemicolon {
		err = fmt.Errorf(`Dummy semicolon is found in import statements: %v`, line)
		return err
	}

	lines = strings.Split(line, ";")

	for _, line := range lines {
		tempLine := strings.TrimSpace(line)
		if len(tempLine) == 0 {
			continue
		}

		doubleQuoteCount := strings.Count(tempLine, `"`)
		singleQuoteCount := strings.Count(tempLine, `'`)

		if doubleQuoteCount == 0 && singleQuoteCount == 0 {
			err = fmt.Errorf(`Missing quote in import statements: %v`, tempLine)
			return err
		}

		if doubleQuoteCount % 2 == 1 || singleQuoteCount % 2 == 1 {
			err = fmt.Errorf(`Mismatch quotes are found in import statements: %v`, tempLine)
			return err
		}

		if doubleQuoteCount > 0 && strings.LastIndex(line, `"`) != (len(tempLine) - 1) {
			err = fmt.Errorf(`Invalid import statements: %v`, tempLine)
			return err
		}

		if singleQuoteCount > 0 && strings.LastIndex(line, `'`) != (len(tempLine) - 1) {
			err = fmt.Errorf(`Invalid import statements: %v`, tempLine)
			return err
		}

		tokens := strings.Split(tempLine, ` `)
		importLib := strings.ReplaceAll(strings.ReplaceAll(tokens[1], `"`, ``), `'`, ``)
		if len(importLib) == 0 {
			err = fmt.Errorf(`Invalid import statements: %v`, tempLine)
			return err
		}

		p.imports = append(p.imports, importLib)

	}

	return nil
}

func (p *Parser) processServiceLines(lines []string) error {
	//nolint: errcheck
	// defer sugar.Sync()

	var err error

	// sugar.Debugf("serviceLines: %v", lines)

	p.services = make([]*Service, 0)

	if len(lines) == 0 {
		return nil
	}

	line := strings.Join(lines, " ")

	serviceCount := strings.Count(line, "service")
	serviceLines := make([]string, 0)
	for i := 0; i < serviceCount; i++ {
		endIndex := strings.Index(line[1:], "service")
		if endIndex == -1 {
			serviceLines = append(serviceLines, line)
		} else {
			serviceLines = append(serviceLines, line[:endIndex])
			line = line[endIndex:]
		}

	}

	for _, serviceLine := range serviceLines {

		if len(serviceLine) == 0 {
			continue
		}

		serviceLine = strings.TrimSpace(serviceLine)

		if serviceLine == ";" {
			err = errors.New("Invalid service declaration")
			return err
		}

		if serviceLine[len(serviceLine) - 1] == ';' {
			err = errors.New("Dummy semicolon at the end of service")
			return err
		}

		beginCurseCount := strings.Count(serviceLine, "{")
		endCurseCount := strings.Count(serviceLine, "}")

		if beginCurseCount != endCurseCount {
			err = fmt.Errorf(`Curse bracket does not match: %v`, serviceLine)
			return err
		}

		beginServiceCursePos := strings.Index(serviceLine, "{")
		endServiceCursePos := strings.LastIndex(serviceLine, "}")

		if beginServiceCursePos == -1 || endServiceCursePos == -1 {
			err = fmt.Errorf(`Cannot find service block: %v`, serviceLine)
			return err
		}

		serviceName := strings.TrimSpace(strings.ReplaceAll(serviceLine[:strings.Index(serviceLine, "{")], "service", ""))
		if strings.Contains(serviceName, " ") || len(serviceName) == 0 {
			err = fmt.Errorf(`Invalid service name: %v`, serviceName)
			return err
		}

		service := &Service{
			name: serviceName,
			rpcs: make([]*Rpc,0),
		}

		rpcBlocks := strings.TrimSpace(serviceLine[strings.Index(serviceLine, "{") + 1:strings.LastIndex(serviceLine, "}")])
		
		if len(rpcBlocks) == 0 {
			err = fmt.Errorf(`No rpc has been found: %v`, serviceLine)
			return err
		}

		rpcBlocks = strings.ReplaceAll(rpcBlocks, "\n", " ")

		for _, rpcBlock := range strings.Split(rpcBlocks, "rpc") {
			rpcBlock = strings.TrimSpace(rpcBlock)

			if len(rpcBlock) == 0 {
				continue
			}

			if rpcBlock[len(rpcBlock) - 1] != ';' {
				err = fmt.Errorf(`Missing semicolon at the end of rpc: %v`, rpcBlock)
				return err
			}

			beginRpcCursePos := strings.Index(rpcBlock, "{")
			endRpcCursePos := strings.LastIndex(rpcBlock, "}")
			if beginRpcCursePos > -1 {
				if endRpcCursePos == -1 {
					err = fmt.Errorf(`Curse bracket mismatch: %v`, rpcBlock)
					return err
				}

				rpcBlock = rpcBlock[:strings.Index(rpcBlock, "{")] + ";"
			} 

			rpcBlock = strings.TrimSpace(strings.ReplaceAll(rpcBlock, ";", ""))

			if len(rpcBlock) == 0 {
				err = fmt.Errorf(`Invalid rpc block: %v`, rpcBlock)
				return err
			}

			firstLeftParenthesisPos := strings.Index(rpcBlock, "(")
			firstRightParenthesisPos := strings.Index(rpcBlock, ")")

			if firstLeftParenthesisPos == -1 {
				err = fmt.Errorf(`Missing Parenthesis: %v`, rpcBlock)
				return err
			}

			if firstRightParenthesisPos == -1 {
				err = fmt.Errorf(`Missing Parenthesis: %v`, rpcBlock)
				return err
			}

			rpcName := rpcBlock[:firstLeftParenthesisPos]
			if len(rpcName) == 0 {
				err = fmt.Errorf(`Cannot obtain rpc name: %v`, rpcBlock)
				return err
			}

			rpcRequest := rpcBlock[firstLeftParenthesisPos + 1:firstRightParenthesisPos]
			if len(rpcRequest) == 0 {
				err = fmt.Errorf(`Cannot find rpc request: %v`, rpcBlock)
				return err
			}

			var rpcResponse string

			returnsPos := strings.Index(rpcBlock, "returns")

			if returnsPos == -1 {
				err = fmt.Errorf(`Missing returns: %v`, rpcBlock)
				return err
			}

			secondLeftParenthesisPos := strings.Index(rpcBlock[returnsPos:], "(")
			secondRightParenthesisPos := strings.Index(rpcBlock[returnsPos:], ")")
			
			if secondLeftParenthesisPos > -1 {
				returnPart := rpcBlock[returnsPos:]
				rpcResponse = returnPart[secondLeftParenthesisPos + 1:secondRightParenthesisPos]
			} else {
				tokens := strings.Split(rpcBlock[returnsPos:], " ")
				rpcResponse = tokens[1]
			}

			if len(rpcResponse) == 0 {
				err = fmt.Errorf(`Cannot find rpc response: %v`, rpcBlock)
				return err
			}

			rpc := &Rpc{
				name: rpcName,
				request: rpcRequest,
				response: rpcResponse,
			}

			fmt.Println("*****BEGIN")
			fmt.Println(rpcBlock)
			fmt.Println(rpcName)
			fmt.Println(rpcRequest)
			fmt.Println(rpcResponse)
			fmt.Println("*****END")

			err = service.addRpc(rpc)
			if err != nil {
				return err
			}

		}

		p.services = append(p.services, service)
	}

	return nil
}

func (p *Parser) processMessageLines(lines []string) error {
	//nolint: errcheck
	// defer sugar.Sync()

	// sugar.Debugf("messageLines: %v", lines)

	return nil
}