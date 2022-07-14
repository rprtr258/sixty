test:
	go test ./... -count=1

todo:
	rg 'TODO' --glob '**/*.go'

# TODO: add linter