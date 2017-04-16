JSX_SRC_DIR=js/src
JS_SRC_DIR=js/src
JSX_BUILD_DIR=build/js
JSX_SRC= main.jsx
JS_SRC= api.js
JSX_OBJS=$(addprefix $(JSX_BUILD_DIR)/,$(JSX_SRC:.jsx=.js))
JS_OBJS=$(addprefix $(JSX_BUILD_DIR)/,$(JS_SRC))
STATIC_BUILD_DIR=build/static
JS_HINT_PATH=node_modules/jshint/bin/jshint --config jshint-config.json
JSCS_PATH=node_modules/jscs/bin/jscs --preset google --config jscs-config.json --fix

all: static js

static: js $(STATIC_BUILD_DIR)/index.html $(STATIC_BUILD_DIR)/style.css $(STATIC_BUILD_DIR)/bootstrap.css $(STATIC_BUILD_DIR)/bootstrap.js

js: build/static/main.js

$(STATIC_BUILD_DIR)/bootstrap.css: bower_components/bootstrap/dist/css/bootstrap.css
	cp $< $@

$(STATIC_BUILD_DIR)/bootstrap.js: bower_components/bootstrap/dist/js/bootstrap.js
	cp $< $@

$(STATIC_BUILD_DIR)/index.html: static/index.html
	cp static/index.html $@

$(STATIC_BUILD_DIR)/style.css: static/style.css
	cp static/style.css $@

$(JSX_BUILD_DIR)/%.js: $(JSX_SRC_DIR)/%.jsx
	$(JSCS_PATH) $<
	babel --presets react $< -o $@
	$(JS_HINT_PATH) $@

$(JSX_BUILD_DIR)/%.js: $(JS_SRC_DIR)/%.js
	$(JSCS_PATH) $<
	cp $< $@
	$(JS_HINT_PATH) $@

build/static/main.js: $(JSX_OBJS) $(JS_OBJS)
	browserify -t debowerify $(JSX_BUILD_DIR)/main.js -o $@

clean:
	rm -rf build
	rm -rf node_modules
	rm -rf bower_components
