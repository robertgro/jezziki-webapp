class Article extends React.Component {
    constructor(props) {
        super(props);
        this.postid = "1"; // default
    }

    render() {
        const postdata = this.props.postdata;
        const indexpost = this.props.indexpost;
        const postext = this.props.postext;


        this.postid = postdata ? postdata.id : indexpost.id;

        const postText = (
            postdata ? (
                postdata.text.indexOf("https") || postdata.text.indexOf("http") != 0 ? (<div dangerouslySetInnerHTML={{__html: postdata.text}} className="p-1"></div>) : (<div dangerouslySetInnerHTML={{__html: "Sie werden weitergeleitet zu ... " + postdata.text}} className="p-1 text-center"></div>)
            ) 
            : (<div dangerouslySetInnerHTML={{__html: indexpost.text}} className="p-1"></div>)
        );

        return (
            <React.Fragment>
                <article className="col-sm-8 p-0 text-wrap" data-ext={postext} id={"p" + this.postid}>
                    {postText}
                </article>
            </React.Fragment>
        );
    }
}