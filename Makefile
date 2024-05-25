.PHONY: push

# 获取当前最新的标签，如果没有则默认为v1.0.0
CURRENT_TAG := $(shell git describe --tags `git rev-list --tags --max-count=1` 2>/dev/null || echo "v1.0.0")

# 生成新的标签（修订号加1）
NEW_TAG := $(shell echo $(CURRENT_TAG) | awk -F. -v OFS=. '{ $$NF = $$NF + 1 ; print }')

# 推送代码并打标签的目标
push:
	# 检查是否有文件更改，如果有则添加并提交它们
	@git diff --exit-code --quiet || (git add . && git commit -m "update")
	# 打上新的标签
	git tag $(NEW_TAG)
	# 推送代码和标签到远程仓库
	git push && git push origin $(NEW_TAG)
	@echo "Pushed changes and tagged as $(NEW_TAG) to origin"
