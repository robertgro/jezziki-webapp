class Aside extends React.Component {
    constructor(props) {
        super(props);
        this.editModeList = [];
    }

    getEditModeElement = element => {
        let elementList = this.editModeList.concat(element);
        this.editModeList = elementList;
    }

    render() {
        const asideitems = this.props.asideitems;
        return (            
            <React.Fragment>
                <aside className="col-sm-2 p-0 text-center">
                    <div className="border-container h-100">
                        <div className="container">
                            <div className="row">
                                <div className="col p-1 text-center">
                                {
                                    asideitems.map((item, index) =>
                                        <button data-em="false" ref={this.getEditModeElement} className="aside-link btn w-75 m-1" onClick={() => {this.props.getpost(item.foreignid, item.ext, this.editModeList[index], false, item.title)}} id={"a" + item.id}>{item.title}</button> 
                                )}
                                </div>
                            </div>
                        </div>
                    </div>
                </aside>
            </React.Fragment>
        );
    }
}