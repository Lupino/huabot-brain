LESS_SOURCE =
JSX_SOURCE = app/dashboard.js \
			 app/dataset.js \
			 app/demo.js \
			 app/searchform.js \
			 app/app.js

HEAD = app/head.js
TAIL = app/tail.js

BUNDLE = app/bundle.js

APP = public/static/js/main.js

all: $(APP)

$(APP): $(BUNDLE)
	browserify -t [ reactify --es6 ] $< | uglifyjs -m -r '$$' > $@

$(BUNDLE): $(HEAD) $(JSX_SOURCE) $(TAIL)
	cat $(HEAD) $(JSX_SOURCE) $(TAIL) > $@


clean:
	rm -f $(APP)
	rm -f $(BUNDLE)
