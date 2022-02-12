class Header extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <React.Fragment>
            <a href="/" target="_self">
                <header className="container bg-dark mw-100">
                    <div className="row h-100">
                        <div className="col-sm">
                            <div className="welcome-textbox mt-2">
                                <h2>Jezziki's Hood</h2>
                            </div>
                            <div className="header-subbox mt-2">
                                <mark>Jezziki.de</mark>
                            </div>
                            <img src="img/logo.png" className="img-fluid" alt="Jezziki.deLogo" />
                        </div>
                    </div>
                </header>
                </a>
            </React.Fragment>
        );
    }
}