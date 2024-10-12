package reviewer

import (
	"code-reviewer/internal/services/llm"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mandolyte/mdtopdf"
	"github.com/openai/openai-go"
)

func Review(filepaths []string) {

	systemmessage := openai.SystemMessage(`
You are an AI that helps people review code. Follow the rules and structure provided below to ensure a thorough and practical review.

#### **Approach to Reviewing Code**
1. **Identify Programming Language**: Recognize the programming language to ensure all feedback aligns with its specific best practices.
2. **Code Quality**: Evaluate readability, adherence to naming conventions, and overall code quality for better maintainability.
3. **Functionality**: Verify that the implemented features fulfill the stated requirements and function correctly.
4. **Testing**: Confirm the presence and adequacy of unit and integration tests, ensuring all edge cases are covered.
5. **Performance**: Identify and recommend areas for performance optimization.
6. **Security**: Identify vulnerabilities like improper input validation, lack of authentication/authorization, and insecure data handling.
7. **Error Handling**: Ensure proper error handling, logging, and prevention of sensitive information leakage.
8. **Code Style and Consistency**: Verify adherence to coding standards and consistency throughout the codebase.
9. **SOLID Principles**:
   - **Single Responsibility**: Each class or function should have only one reason to change.
   - **Open/Closed**: Code should be open for extension but closed for modification.
   - **Liskov Substitution**: Derived classes should be usable in place of their base classes.
   - **Interface Segregation**: Use multiple small interfaces rather than a single large one.
   - **Dependency Inversion**: Depend on abstractions, not concrete implementations.

#### **Key Concerns During Code Review**
1. **Code Clarity and Readability**: Check for descriptive variable and function names, and avoid overly complex code.
2. **Maintainability**: Ensure ease of modification and extensibility by identifying tight coupling and promoting modular design.
3. **Duplication**: Suggest refactoring to eliminate repeated code and promote reusability.
4. **Scalability and Performance**: Highlight potential performance bottlenecks, such as inefficient loops or redundant computations.
5. **Security**: Identify common vulnerabilities like SQL injection, XSS, and the use of hardcoded secrets.
6. **Testing Gaps**: Ensure that there is sufficient test coverage, including edge cases.
7. **Logical Errors and Bugs**: Identify logical issues or bugs in the code.
8. **Error Handling and Edge Cases**: Confirm comprehensive error handling and suitable behavior for edge cases.
9. **Consistency**: Ensure consistent use of code style and design patterns throughout.
10 **Cognitive Complexity** : Ensure code easy to manage, easy to understand and easy to maintenance


#### **Review Guidelines**
- **Limit Scope**: Focus only on the provided code; ignore external dependencies or functions.
- **No Assumptions**: Assume external APIs and functions work as expected; do not comment on their validity.

#### **Review Report Format**
1. **General Feedback**: Provide a high-level overview of the code quality.
2. **Detailed Review**:
   - **Line Number**: Specify the line of code in question.
   - **Issue**: Clearly describe the problem.
   - **Original Code Snippet**: Include relevant code for context.
   - **Suggestion**: Provide an improved snippet  code or actionable advice, if the code need to be breakdown or refactor give the full code.
3. **Strengths and Summary**: Highlight well-written parts of the code, and summarize the key findings of the review.

**Example Review Format**:
1. **Line 42: Error Handling**
   - **Issue**: The error from os.Open is not properly handled.
   - **Original Code**:
	` + "```go" + `
     file, err := os.Open(path)
     if err != nil {
         fmt.Errorf("Failed to open file", err)
     }
     ` + "```" + `
   - **Suggestion**:
     The error should be logged or cause program termination for better handling:
     ` + "```go" + `
     file, err := os.Open(path)
     if err != nil {
         log.Fatalf("Failed to open file: %v", err)
     }
     ` + "```" + `

#### **General Principles**
- **Focus Only on Visible Issues**: Comment only on the issues explicitly present in the provided code.
- **Constructive Feedback**: Provide clear, practical suggestions while avoiding unnecessary complexity.
- **Practical Improvements**: Align recommendations with the code's complexity and intended purpose.

#### **Reporting Rules**
- **Do Not Provide Full Code if only minor suggestion, you able to provide full code if the code need to breakdown or refactor**
- **Detailed Snippets**: Provide sufficiently detailed code snippets to avoid confusion, avoiding overly brief suggestions.
- **Markdown Structure**: Use the following headers to organize the review:
  ### [Chapter Number].1 General Feedback - [Filename]
  ...
  ### [Chapter Number].2 Detailed Review - [Filename] 
  ...
  ###  [Chapter Number].3 Summary - [Filename] 
  ...
- DO NOT USE ` + "`" + ` in header and for filename, do not put fullpath,just filename only

### **Markdown Rules**
- **Use best practise markdown format
- **Avoid MD032/blanks-around-lists**: Lists should be surrounded by blank 
- **Avoid MD022/blanks-around-headings**: Headings should be surrounded by blank lines
- **Avoid MD009/no-trailing-spaces**: Trailing spaces 
- **Avoid MD031/blanks-around-fences**: Fenced code blocks should be surrounded by blank 


IF you the script is good, you can skip it, put #NA only for response

`)
	var results []string
	mdfilename := "code_review.md"
	mdfile, err := os.Create(mdfilename)
	if err != nil {
		log.Printf("Failed to create markdown file: %v", err)
	}
	defer mdfile.Close()

	for index, path := range filepaths {
		func() {
			log.Printf("Reviewing filepath : %s\n", path)
			file, err := os.Open(path)
			if err != nil {
				log.Printf("Failed to open file: %v", err)
				return
			}
			defer file.Close()

			content, err := io.ReadAll(file)
			if err != nil {
				log.Printf("Failed to read file %s: %v", path, err)
				return
			}

			var messages []openai.ChatCompletionMessageParamUnion
			messages = append(messages, systemmessage)
			stringContent := "File Path: " + path + "\n"
			stringContent += "Chapter Number : " + strconv.Itoa(index+1) + "\n\n"
			stringContent += string(content)

			messages = append(messages, openai.UserMessage(stringContent))
			tempresult := "##" + strconv.Itoa(index+1) + " File: " + path
			// res, err := runChat(messages)
			res, err := llm.ChatMPN1(messages)
			if err != nil {
				log.Printf("Error while running chat concurent %v", err)
			}
			tempresult += "\n\n" + res
			results = append(results, tempresult)

			_, err = mdfile.WriteString(tempresult + "\n\n")
			if err != nil {
				log.Printf("Failed to write result to file: %v", err)
				return
			}

			log.Printf("Reviewing filepath : %s..... Done\n", path)
		}()
	}

	// convertMDtoPDF(mdfilename, "code_review.pdf")
	convertMDtoHTML(mdfilename, "code_review.html")

}

func runChat(messages []openai.ChatCompletionMessageParamUnion) (string, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	// Channels to receive results and errors from the goroutines
	resChan := make(chan string, 2)
	errChan := make(chan error, 2)

	// Run llm.ChatMPN1 concurrently
	go func() {
		defer wg.Done()
		res, err := llm.ChatMPN1(messages)
		resChan <- res
		errChan <- err
	}()

	// Run llm.ChatMPN2 concurrently
	go func() {
		defer wg.Done()
		res2, err2 := llm.ChatMPN2(messages)
		resChan <- res2
		errChan <- err2
	}()

	// Wait for both goroutines to finish
	wg.Wait()
	close(resChan)
	close(errChan)

	// Retrieve results and errors
	var res, res2 string
	var err, err2 error

	for r := range resChan {
		if res == "" {
			res = r
		} else {
			res2 = r
		}
	}

	for e := range errChan {
		if err == nil {
			err = e
		} else {
			err2 = e
		}
	}

	if err != nil {
		log.Printf("Error when getting review from LLM for MPN1 file : %v", err)
	}

	if err2 != nil {
		log.Printf("Error when getting review from LLM for MPN2 file : %v", err2)
	}
	_ = res2

	// log.Println("")
	// log.Printf("Response1 : %s\n", res)

	// log.Println("")

	return res, err
}

func convertMDtoPDF(mdfile string, pdfile string) {

	inputFile, err := os.ReadFile(mdfile)
	if err != nil {
		log.Fatalf("Failed to Read MD File: %v", err)
	}

	var opts []mdtopdf.RenderOption

	renderer := mdtopdf.NewPdfRenderer("", "", pdfile, "", opts, mdtopdf.LIGHT)
	if err := renderer.Process(inputFile); err != nil {
		log.Fatalf("Failed to render PDF: %v", err)
	}

}

func convertMDtoHTML(mdfile string, htmlfile string) {
	mdContent, err := os.ReadFile(mdfile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	// Create a Markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(mdContent)

	// Create an HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Convert Markdown to HTML
	htmlContent := markdown.Render(doc, renderer)

	// Write the HTML content to the output file
	err = os.WriteFile(htmlfile, htmlContent, 0644)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("Converted %s to %s\n", mdfile, htmlfile)
}
