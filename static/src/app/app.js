class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      postdata: null,
      postext: null,
      navitems: [],
      asideitems: [],
      footeritems: [],
      indexpost: null,
    };
    this.getPostbyID = this.getPostbyID.bind(this);
    this.getIndex = this.getIndex.bind(this);
    this.loadExternal = this.loadExternal.bind(this);
    this.applyWindowTop = this.applyWindowTop.bind(this);
    this.changeDocTitle = this.changeDocTitle.bind(this);
  }

  changeDocTitle(title) {
    var docTitle = document.title;
    var domainTitle;
    if(docTitle.includes(" - ")) {
      domainTitle = docTitle.slice(0, docTitle.indexOf(" "));
      document.title = domainTitle + " - " + title;
    } else {
      document.title = docTitle + " - " + title;
    }
  }

  getPostbyID(id, ext, em, wm, title) {
      this.changeDocTitle(title);
      fetch("/api/v1/posts/" + id)
      .then(response => {
        if (!response.ok) {
          return response.text().then(text => {throw new Error(text)})
        } 
        return response.json();
      })
      .then(data => {
        this.setState({ 
          postdata: data, 
          postext: ext, 
        });
        this.loadExternal(ext, em);
      })
      .catch(err => {
        console.log(err);
      });
      this.applyWindowTop(em, wm);
  }

  applyWindowTop(em, wm) {
    if(em.getAttribute("data-em") === "false" && Boolean(wm)) {
      window.scrollTo(0,0);
    }
  }

  loadExternal(ext, em) {
    if(ext === "true" && em.getAttribute("data-em") === "false") {
      var url = this.state.postdata.text.replace("/\/$/","").replace(/\u200B/g,"").replace("\ufeff", "");
      // UTF-8 BOM removal https://stackoverflow.com/questions/13024978/removing-bom-characters-from-ajax-posted-string
      // window open adds a trailing slash https://stackoverflow.com/questions/61769019/window-open-adds-exrtra-forward-slash-before-the-query-parameters-start
      // remove zero width space https://stackoverflow.com/questions/24205193/javascript-remove-zero-width-space-unicode-8203-from-string
      window.open(url, '_blank').focus();
    }
  }

  getIndex() {
    fetch("/api/v1/index")
    .then(response => {
      if (!response.ok) {
        return response.text().then(text => {throw new Error(text)})
      } 
      return response.json();
    })
    .then(data => {
      this.setState({
        navitems: data.nav,
        asideitems: data.aside,
        footeritems: data.footer,
        indexpost: data.indexpost,
        postext: false,
      });
    })
    .catch(err => {
      console.log(err);
    });
  }

  componentDidMount() {
    this.getIndex();
  }

  render() {
    return (
      <React.Fragment>
      { this.state.indexpost && 
        <React.Fragment>
              <Header />
              <Main navitems={this.state.navitems} asideitems={this.state.asideitems} postdata={this.state.postdata} postext={this.state.postext} indexpost={this.state.indexpost} getpost={this.getPostbyID} />
              <Footer footeritems={this.state.footeritems} getpost={this.getPostbyID} />
        </React.Fragment>
      }
      </React.Fragment>
    );
  }
}