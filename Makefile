TAG ?= v0.0.1
MSG ?= "release: new version"

.PHONY: add commit push tag release

add:
	git add .

commit:
	git commit -m '$(MSG)'

push:
	git push

tag:
	git tag $(TAG)
	git push origin $(TAG)

release: add commit push tag
