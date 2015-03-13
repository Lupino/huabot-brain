LESS_SOURCE =
JSX_SOURCE = app/dashboard.js \
			 app/dataset.js \
			 app/demo.js \
			 app/searchform.js \
			 app/app.js

HEAD = app/head.js
TAIL = app/tail.js

APP = public/static/js/main.js

all: $(APP)

$(APP): comibed.js
	browserify -t [ reactify --es6 ] comibed.js | uglifyjs -m -r '$$' > $@

comibed.js: $(HEAD) $(JSX_SOURCE) $(TAIL)
	cat $(HEAD) $(JSX_SOURCE) $(TAIL) > comibed.js


clean:
	rm -f $(APP)
	rm -f comibed.js
