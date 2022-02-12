class Main extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        const postdata = this.props.postdata;
        const postext = this.props.postext;
        const indexpost = this.props.indexpost;
        const asideitems = this.props.asideitems;
        const navitems = this.props.navitems;
        return (
            <div>
                <main className="container bg-dark mw-100">
                    <section className="row">
                        <Nav navitems={navitems} getpost={this.props.getpost} />
                        <Article postdata={postdata} indexpost={indexpost} postext={postext} />
                        <Aside getpost={this.props.getpost} asideitems={asideitems} postdata={postdata} />
                    </section>
                </main>
            </div>
        );
    }
}