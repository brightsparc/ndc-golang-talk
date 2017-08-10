# ndc-golang-talk

## From C# to Golang, My data science journey

This repository contains my presentation and code examples for my [NDC Talk](http://ndcsydney.com/talk/from-c-to-golang-my-data-science-journey./)

### Run Presentation

Run presentation with [Go Presentation tool](https://godoc.org/golang.org/x/tools/present) for my slides.  Get started with this [Tutorial](http://halyph.com/blog/2015/05/18/golang-presentation-tool.html) for getting setup on a mac.

```
go install golang.org/x/tools/cmd/present
```

Optionally install [Syntax Highlighting](https://github.com/JosephBuchma/Go-Present-code-highlighter)

### Demonstration

Source code for demonstrate is within this repository:

1. [NDC Segment](https://github.com/brightsparc/ndc-golang-talk/tree/master/ndc_segment) Go API for event tracking.
2. [NDC Scraper](https://github.com/brightsparc/ndc-golang-talk/tree/master/ndc_scraper) Python scripts to scrape NDC speaker talks
3. [NDC Predict](https://github.com/brightsparc/ndc-golang-talk/tree/master/ndc_predict) Go API predicting talk tags based on content.
