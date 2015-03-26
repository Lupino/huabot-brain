LESS_SOURCE =
JSX_SOURCE = app/dashboard.jsx \
			 app/dataset.jsx \
			 app/demo.jsx \
			 app/searchform.jsx \
			 app/app.jsx

HEAD = app/head.jsx
TAIL = app/tail.jsx

BUNDLE = app/bundle.jsx

APP = public/static/js/main.js

all: $(APP)

$(APP): $(BUNDLE)
	browserify -t [ reactify --es6 envify --NODE_ENV production] $< | uglifyjs -m -r '$$' > $@

$(BUNDLE): $(HEAD) $(JSX_SOURCE) $(TAIL)
	cat $(HEAD) $(JSX_SOURCE) $(TAIL) > $@

deps: package.json
	npm install

clean:
	rm -f $(APP)
	rm -f $(BUNDLE)
