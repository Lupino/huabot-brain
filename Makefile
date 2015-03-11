LESS_SOURCE =
JSX_SOURCE = app/head.js \
			 app/dashboard.js \
			 app/dataset.js \
			 app/demo.js \
			 app/searchform.js \
			 app/app.js \
			 app/tail.js
APP = public/static/js/main.js

all: $(APP)

$(APP): $(JSX_SOURCE)
	cat $(JSX_SOURCE) > comibed.js
	browserify -t [ reactify --es6 ] comibed.js | uglifyjs -m -r '$$' > $@


clean:
	rm -f $(APP)
	rm -f comibed.js
