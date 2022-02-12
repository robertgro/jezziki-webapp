class Footer extends React.Component {
    constructor(props) {
        super(props);
        this.editModeList = [];
    }

    getEditModeElement = element => {
        let elementList = this.editModeList.concat(element);
        this.editModeList = elementList;
    }

    render() {
        const footeritems = this.props.footeritems;
        return (
            <React.Fragment>
                <footer className="container mw-100">
                    <div className="row h-100">
                        <div className="col-sm">
                            <div className="d-flex flex-row flex-wrap justify-content-center mt-2">
                            {
                                footeritems.map((item, index) =>
                                <a data-em="false" ref={this.getEditModeElement} role="button" onClick={() => { this.props.getpost(item.foreignid, item.ext, this.editModeList[index], true, item.title)}} id={"f"+item.id}>{item.title}</a>
                            )}
                            </div>
                        </div>
                    </div>
                </footer>
            </React.Fragment>
        );
    }
}