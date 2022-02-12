class Root extends React.Component {
  constructor(props) {
    super(props);
  }
  render() {
    return (
      <React.Fragment>
        <App />
      </React.Fragment>
    );
  }
}
ReactDOM.render(
  <Root />,
  document.getElementById('root')
);

