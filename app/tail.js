var routes = (
  <Route handler={App} path="/">
    <DefaultRoute handler={Dashboard} />
    <Route name="dashboard" handler={Dashboard} />
    <Route name="demo" handler={DEMO} />
    <Route name="datasets" handler={Datasets} path="/ds/:dataType"/>
    <NotFoundRoute handler={Datasets} />
  </Route>
);

ReactRouter.run(routes, function (Handler, state) {
  React.render(<Handler />, document.body);
});
