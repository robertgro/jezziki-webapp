class Nav extends React.Component {
    constructor(props) {
        super(props);
        this.renderNavTree = this.renderNavTree.bind(this);
        this.parentnavid = "";
        this.editModeList = [];
    }

    getEditModeElement = element => {
        let elementList = this.editModeList.concat(element);
        this.editModeList = elementList;
    }

    setParentNavID(id) {
        this.parentnavid = id;
    }

    // credits to recursive tree render https://stackoverflow.com/questions/45790499/react-json-tree-with-the-help-of-recursion
    renderNavTree = (node) => this.props.navitems.filter(itemNode => itemNode.parentid == node).map((childNode, index) =>
                    (<div className="nav-div mx-auto" id={"n" + childNode.id}>
                    {childNode.parentid === 0 
                    ? (<ul className="navbar-nav flex-column">{this.setParentNavID("node_" + childNode.id)}<li className="nav-item"><a className="nav-link" data-em="false" ref={this.getEditModeElement} role="button" data-toggle="collapse" data-target={"#" + this.parentnavid} onClick={() => this.props.getpost(childNode.foreignid, childNode.ext, this.editModeList[index], false, childNode.title)}>{childNode.title}</a></li>{this.renderNavTree(childNode.id)}</ul>) 
                    : (<ul className="navbar-nav flex-column collapse" id={this.parentnavid}><li className="nav-item-c"><a className="nav-link" data-em="false" ref={this.getEditModeElement} role="button" onClick={() => this.props.getpost(childNode.foreignid, childNode.ext, this.editModeList[index], false, childNode.title)}>{childNode.title}</a></li>{this.renderNavTree(childNode.id)}</ul>)
                    }
                    </div>))

    render() {
        return (
            <React.Fragment>
                <nav className="p-0 col-sm-2 navbar navbar-expand-xl bg-dark navbar-dark d-block" role="navigation">
                    <div className="inner-border-container mt-2 mb-2">
                        <div className="text-center">
                            <button className="mt-1 mb-2 navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarcontent">
                                <span class="navbar-toggler-icon"></span>
                            </button>
                        </div>
                        <div className="navbar-collapse collapse flex-column" id="navbarcontent">
                            {this.renderNavTree(0)}
                        </div>
                    </div>
                </nav>
            </React.Fragment>
        );
    }
}