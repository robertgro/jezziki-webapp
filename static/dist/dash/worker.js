
class editContentWorker {

    ContentMain = document.getElementsByTagName("article")[0].children[0]; // div inside cmain

    constructor() {
		this.ContentMain.insertAdjacentElement('afterbegin', this.CreateContentEditButton());
    }
    
    CreateContentEditButton() {
        var EditButton = document.createElement("button");
        EditButton.onclick = this.editContent.bind(this);
        EditButton.setAttribute("class","btn btn-danger float-right m-2")
        EditButton.setAttribute("id", "editContentBtn");
        EditButton.textContent = "Edit"
        return EditButton
    }
    
    editContent() {
        document.getElementById("editContentBtn").remove();
        document.getElementById("editNavBtn").disabled = true;
        document.getElementById("editAsideBtn").disabled = true;
        document.getElementById("editFooterBtn").disabled = true;
        
        let ContentMainIH = this.ContentMain.innerHTML;
        this.ContentMain.innerHTML = "";
    
        let NewTextareaObject = this.prepareContentEditDOM();
        let NewSaveContentBtn = this.CreateContentSaveButton();
    
        this.loadExternalSwitch();
        this.ContentMain.appendChild(NewTextareaObject);
        this.ContentMain.appendChild(NewSaveContentBtn);
        this.newJoditEditor(ContentMainIH);
        document.getElementsByClassName('jodit-container')[0].classList.add("m-2");
        document.getElementsByClassName("jodit-placeholder")[0].remove();
    }

    loadExternalSwitch() {
        let checkedSwitch = "";
        if(this.ContentMain.parentElement.getAttribute("data-ext") === "true") {
            checkedSwitch = "checked";
        }
        let switchHtml = '<div class="form-check m-2 text-center"> \
        <input class="form-check-input" onclick="ContentonClickAddLabel()" type="checkbox" value="" id="extSwitch" ' + checkedSwitch + '> \
        <label class="form-check-label" for="extSwitch"> \
          External \
        </label> \
      </div>';
        this.ContentMain.insertAdjacentHTML('afterbegin', switchHtml);
    }
    
    prepareContentEditDOM() {
        var TextareaObject = document.createElement("textarea");
        TextareaObject.setAttribute("id","editor");
        TextareaObject.setAttribute("name","editor");
        //console.log(TextareaObject);
        return TextareaObject
    }
    
    // credits https://rmamuzic.rs/node_modules/jodit/examples/custom-toolbar.html
    newJoditEditor(contentInnerHtml) {
        var customButtonBar = [
            'source',
            '|',
            'bold',
            'italic',
            'underline',
            'strikethrough',
            'hr',
            'symbol',
            'eraser',
            '|',
            'ul',
            'ol',
            'paragraph',
            'align',
            '|',
            'image',
            'video',
            'link',
            'table',
            'superscript',
            'subscript',
            '|',
            'paste',
            'selectall',
            'undo',
            'redo',
            'find',
            '|',
            'fullsize',
            'preview',
            'about',
        ]
        
        var editor = new Jodit('#editor',{
            buttons: customButtonBar,
            buttonsMD: customButtonBar,
            buttonsSM:  customButtonBar,
            buttonsXS: customButtonBar,
            height: 500,
            theme: "dark",
        });

        editor.value = contentInnerHtml; 
    }
    
    CreateContentSaveButton() {
        var SaveButton = document.createElement("button");
        SaveButton.onclick = this.saveContent.bind(this);
        SaveButton.setAttribute("class","btn btn-success float-right m-2")
        SaveButton.setAttribute("id", "saveContentBtn");
        SaveButton.textContent = "Save"
        return SaveButton
    }

    parseContentChildElements(element) {
        if(element.length > 0) {
            for(var i = 0; i < element.length; i++) {
                var e = element[i];

                switch(e.tagName) {
                    case "IFRAME":
                        var src = e.getAttribute("src");
                        var x;
                        if(src.includes("www.example.com")) {
                            x = src.replace("www.example.com","localhost");
                        } else if(src === "undefined") {
                            x = "https://player.twitch.tv/?channel=jezziki&parent=localhost";
                        } else {
                            x = src;
                        }
                        e.setAttribute("src", x);
                        if(e.getAttribute("width").replace("px","") <= 350) {
                            if(!e.classList.contains("float-left")) {
                                e.classList.add("float-left");
                            }
                        } else {
                            if(e.classList.contains("float-left")) {
                                e.classList.remove("float-left");
                            }
                        }
                        if(!e.classList.contains("m-2")) {
                            e.classList.add("m-2");
                        }
                        break;
                    case "IMG":
                         //console.log(e.getAttribute("width"), e.classList);
                        if(!e.hasAttribute("width")) {
                            e.setAttribute("width", "300");
                        }

                        if(e.getAttribute("width").replace("px","") <= 350) {
                            if(!e.classList.contains("float-left")) {
                                e.classList.add("float-left");
                            }
                        } else {
                            if(e.classList.contains("float-left")) {
                                e.classList.remove("float-left");
                            }
                        }
                        if(!e.classList.contains("m-2")) {
                            e.classList.add("m-2");
                        }
                        break;
                    default:
                        break;
                }
            }
        }
    }
    
    // recursively , greetings to https://stackoverflow.com/questions/2161634/how-to-check-if-element-has-any-children-in-javascript
    parseContentStyleAttribute(element) {
        if(element.length > 0) {
            for(var i = 0; i < element.length; i++) {
                var e = element[i];
                
                if(e.hasAttribute("style")) {
                    parseStyleAttribute(e);
                }
                
                e.removeAttribute("style");

                if(e.childElementCount > 0) {
                    this.parseContentStyleAttribute(e.children);
                }
            }
        }
    }
    
    saveContent() {
        var JoditEditor = document.getElementsByClassName('jodit-wysiwyg')[0];
        var JoditTags = document.querySelectorAll('jodit');

        // remove empty p tags
        document.querySelectorAll('.jodit-wysiwyg > *').forEach(function(child){
            if(child.childElementCount === 1 && child.children[0].tagName === 'BR' && child.tagName === 'P') {
                child.remove();
            }
        });
    
        // if any embbedded jodit tag, unfold
        if(JoditTags.length > 0) {
            //console.log(JoditTags.length);
            for(var i = 0; i < JoditTags.length; i++) {
                var JoditTag;
                JoditTag = JoditTags[i];
                //console.log(JoditTag, JoditTags[i], i);
                unwrap(JoditTag);
            }
        }

        // edit iframe, image attributes/classes, remove empty paragraphs
        this.parseContentChildElements(JoditEditor.children);
    
        // edit style attribute due to csp recursively
        this.parseContentStyleAttribute(JoditEditor.children);

        var ext = document.getElementById('extSwitch').checked;
    
        // remove everything to be replaced with the editors content afterwards
        while(this.ContentMain.hasChildNodes()) {
            this.ContentMain.removeChild(this.ContentMain.firstChild);
        }
    
        this.ContentMain.innerHTML = JoditEditor.innerHTML;
    
        this.createPostPayload(this.ContentMain.parentNode.id.slice(1), this.ContentMain.innerHTML.replace(/\u200B/g,''), ext);
    
        document.getElementById("editNavBtn").disabled = false;
        document.getElementById("editAsideBtn").disabled = false;
        document.getElementById("editFooterBtn").disabled = false;
    }
    
    createPostPayload(id, content, ext) {
        var jsonObject = {"id":id, "text": content.trim(), "ext": ext.toString()};
        var jsonExt;
        var jsonFin;
        //console.log(JSON.stringify(jsonObject));
        // Fix for external posts
        if(jsonObject.text.includes("<p>Sie werden weitergeleitet zu ...")) {
            jsonExt = jsonObject.text.replace("</p>","").replace("</a>","");
            if(jsonExt.lastIndexOf("https") != -1) {
                jsonFin = jsonExt.slice(jsonExt.lastIndexOf("https"));
            } else if(jsonExt.lastIndexOf("http") != -1) {
                jsonFin = jsonExt.slice(jsonExt.lastIndexOf("https"));
            } else {
                console.log("Please provide a link with a protocol e.g. https://example.com http://secexample.net/");
                alert("Bitte einen Link mit HTTPS:// oder HTTP:// vorher eingeben! Setze Link auf: https://twitch.tv/jezziki");
                jsonFin='https://twitch.tv/jezziki';
            }
            jsonObject.text = jsonFin; // https://stackoverflow.com/questions/1067742/how-can-i-clean-source-code-files-of-invisible-characters
            // https://stackoverflow.com/questions/24205193/javascript-remove-zero-width-space-unicode-8203-from-string
            //console.log(jsonExt, jsonObject.text, jsonFin);
        }
        console.log(jsonFin, jsonObject.text);
        updateComponent("/5fzt78g4A7fnb882/post", JSON.stringify(jsonObject), false);
    }

}

class editNavWorker {

    NavMain = document.getElementsByTagName("nav")[0];
    NavItem = {};
    NavItemList = [];

    constructor() {
        this.NavMain.insertAdjacentElement('beforeend', this.createNavEditButton());
        this.getPreviousNav();
    }

    getPreviousNav() {
        let navItems = document.querySelectorAll("nav li");
        this.NavItemList = [];

        for(let item of navItems) {
            this.NavItem = {};
            this.NavItem.id = item.parentElement.parentElement.id.slice(1); // <li> <ul> <div>
            if(item.parentElement.id === "") {
                this.NavItem.parentid = "0";
            } else {
                this.NavItem.parentid = item.parentElement.id.slice(5); // cut 'node_'
            }
            this.NavItem.title = item.innerHTML;
            this.NavItemList.push(this.NavItem);
        }
    }

    createNavEditButton() {
        var EditButton = document.createElement("button");
        EditButton.onclick = this.editNav.bind(this);
        EditButton.setAttribute("class","btn btn-danger float-left ml-2")
        EditButton.setAttribute("id", "editNavBtn");
        EditButton.textContent = "Edit"
        return EditButton
    }

    createNavSaveButton() {
        var SaveButton = document.createElement("button");
        SaveButton.onclick = this.saveNav.bind(this);
        SaveButton.setAttribute("class","btn btn-success float-right mr-2")
        SaveButton.setAttribute("id", "saveNavBtn");
        SaveButton.textContent = "Save"
        return SaveButton
    }

    saveNav() {
        document.getElementById("saveNavBtn").remove();
        document.getElementById("editContentBtn").disabled = false;
        document.getElementById("editAsideBtn").disabled = false;
        document.getElementById("editFooterBtn").disabled = false;
        document.querySelectorAll("ul i").forEach((item) => item.remove());
        document.querySelectorAll("nav a").forEach((item) => item.setAttribute("data-em", "false"))
        this.NavMain.insertAdjacentElement('beforeend', this.createNavEditButton());
        document.querySelectorAll('nav ul').forEach((item) => item.querySelectorAll("li").forEach((item) => item.ondblclick = null));
        this.getNextNav();
        this.createNavPayload();
    }

    createNavPayload() {
        var jsonObject = {"navItems":this.NavItemList};
        console.log(JSON.stringify(jsonObject));
        updateComponent("/5fzt78g4A7fnb882/nav", JSON.stringify(jsonObject), true);
        this.NavItemList = [];
    }

    getNavItemIndex() {
        for(let j=0; j < this.NavItemList.length; j++) {
            if(this.NavItemList[j].ID === this.NavItem.id) {
                return j;
            }
        }
        return -1;
    }

    getNextNav() {
        let NextNavList = document.querySelectorAll("nav li");
        let i;
        this.NavItemList = [];
        //NextNavList.forEach((item, index) => console.log(index, item.parentElement.id.slice(1), item.innerHTML));

        for(i=0; i < NextNavList.length; i++) {
            this.NavItem = {};
            this.NavItem.id = NextNavList[i].parentElement.parentElement.id.slice(1);
            this.NavItem.title = NextNavList[i].children[0].innerHTML; // <li> -> <a>[0]

            if(NextNavList[i].parentElement.id === "") {
                this.NavItem.parentid = "0";
            } else {
                this.NavItem.parentid = NextNavList[i].parentElement.id.slice(5); // cut 'node_'
            }

            var index = this.getNavItemIndex();
            
            if(index === -1) {
                this.NavItemList.push(this.NavItem);
            } else {
                this.NavItemList[index] = this.NavItem;
            }
        }        
        // this.NavItemList.forEach((item, index) => console.log(index, item.ID, item.Title));
    }

    applyEditButtons() {
        var navItems = document.querySelectorAll('nav ul');
        var editBtns = document.querySelectorAll('nav i');
        editBtns.forEach((item) => item.remove()); // remove if any
        
        for(let item of navItems) {

            // remove if any
            item.querySelectorAll("li").forEach((item) => item.removeEventListener('dblclick', this.editNavCallback));

            // first item can't be deleted to retain nav
            if(!(item.parentElement.id.slice(1) === '1')) {
                item.insertAdjacentElement('afterbegin', this.createDelNavItemButton());
            }

            item.insertAdjacentElement('afterbegin', this.createAddNavItemButton());

            // For subitem: <ul> -> <div> -> <ul> -> <DIV id="parent"> | Else: <ul> -> <div> -> <div> -> <div> 
            let parentElement = item.parentElement.parentElement.parentElement;

            // no subcat for subitems (without parentID) & has no subitems yet (count less than 4)
            if((parentElement.id === "" && item.children.length < 4)) {
                // and is not (first nav item and has more than two children) due to first nav item has no delete button
                if(!(item.parentElement.id.slice(1) === '1' && item.children.length > 2)) {
                    item.insertAdjacentElement('afterbegin', this.createAddSubNavItemButton());
                }
            }

            item.querySelectorAll("li").forEach((item) => item.ondblclick = this.editNavItemCallback.bind(this));
        }
    }

    // for prod for(item of document.getElementsByClassName("nav-div")) {console.log(item);}
    editNav() {
        document.getElementById("editContentBtn").disabled = true;
        document.getElementById("editAsideBtn").disabled = true;
        document.getElementById("editFooterBtn").disabled = true;
        document.getElementById("editNavBtn").remove();
        document.querySelectorAll("nav a").forEach((item) => item.setAttribute("data-em", "true"))
        this.NavMain.insertAdjacentElement('beforeend', this.createNavSaveButton());
        this.applyEditButtons();
    }

    editNavItemCallback(e) {
        e.stopPropagation();
        let formerValue = e.target.innerHTML;
        let editInput = document.createElement('input');
        editInput.type = "text";
        editInput.id = "editInput";
        editInput.size = "11";
        editInput.maxLength = "25";
        editInput.value = formerValue;
        editInput.onfocusout = this.saveNavInput.bind(this);
        e.target.innerHTML = "";
        e.target.appendChild(editInput);
        document.getElementById('editInput').focus();
    }

    saveNavInput(e) {
        let inputTag = document.getElementById('editInput');
        if(this.NavMain.contains(inputTag)) {
            let inputValue = inputTag.value;
            inputTag.parentNode.innerHTML = inputValue;
        }
    }

    createAddNavItemButton() {
        var AddButton = document.createElement("i");
        AddButton.setAttribute("class","bi bi-plus-square navaddbtn")
        AddButton.addEventListener('click', this.addNavItem.bind(this));
        return AddButton
    }

    createDelNavItemButton() {
        var DelButton = document.createElement("i");
        DelButton.setAttribute("class","bi bi-dash-square navdelbtn")
        DelButton.addEventListener('click', this.delNavItem.bind(this));
        return DelButton
    }

    createAddSubNavItemButton() {
        var AddSubButton = document.createElement("i");
        AddSubButton.setAttribute("class","bi bi-box-arrow-down navaddsubbtn")
        AddSubButton.addEventListener('click', this.addSubNavItem.bind(this));
        return AddSubButton
    }

    getNavIndexList() {
        let navItems = document.querySelectorAll("nav ul");
        let navIndexList = []; // need an array here for splice (not applicable to navItems htmlcollection)
        navItems.forEach((_, index) => navIndexList.push(index)); // get index  from current nav item
        return navIndexList;
    }

    // insert a new nav item, assign edit buttons and reindex, e.target is <i> here
    addNavItem(e) {
        let navIndexList = this.getNavIndexList();
        let parentDiv = e.target.parentElement.parentElement; // <i> -> <ul> -> <div>
        let newNavHtml;
        let parentID = parentDiv.parentElement.parentElement.id; // <div> -> <ul> -> <div id="parent"> if any
        let nodeID = (parseInt(parentDiv.id.slice(1)) + 1); // calculate nodeID

        if(parentID === "") {
            newNavHtml = '<div class="nav-div mx-auto"><ul class="navbar-nav flex-column"><li class="nav-item"><a class="nav-link" data-em="false" role="button" data-toggle="collapse" data-target="#node_' + nodeID + '">Neues Item</a></li></ul></div>';
            e.target.parentElement.parentElement.insertAdjacentHTML('afterend', newNavHtml);
            
        } else {
            newNavHtml = '<div class="nav-div mx-auto"><ul class="navbar-nav flex-column collapse show" id="' + parentID + '"><li class="nav-item-c"><a class="nav-link" data-em="false" role="button">Neues Subitem</a></li></ul></div>';
            e.target.parentElement.parentElement.insertAdjacentHTML('afterend', newNavHtml); // insert new html
        }

        navIndexList.splice(0, 0, 0); // insert new index
        this.ReIndexNavItems(navIndexList); // assign id's to items
        this.applyEditButtons(); // remove and re-assign edit buttons
    }

    ReIndexNavItems(navIndexList) {
        let navItems = document.querySelectorAll("nav ul");
        for(let i = 0; i < navIndexList.length; i++) {
            let index = i + 1;
            navItems[i].parentElement.id = "n" + index; // ids beginning with 1

            let subItems = navItems[i].querySelectorAll("ul");

            if(subItems.length > 0) {
                navItems[i].querySelector("a").setAttribute("data-target", "#node_" + index);
                subItems.forEach((item) => item.id = "node_" + index);
            }
        }
    }

    addSubNavItem(e) {
        let navIndexList = this.getNavIndexList();
        let newNavHtml = '<div class="nav-div mx-auto"><ul class="navbar-nav flex-column collapse show" id="' + e.target.parentElement.parentElement.id + '"><li class="nav-item-c"><a class="nav-link" data-em="false" role="button">Neues Subitem</a></li></ul></div>';
        navIndexList.splice(0, 0, 0); // insert new index
        e.target.parentElement.insertAdjacentHTML('beforeend', newNavHtml); // insert new html
        this.ReIndexNavItems(navIndexList); // assign id's to items
        this.applyEditButtons(); // remove and re-assign edit buttons
    }

    delNavItem(e) {
        let ListTag = e.target.parentElement; // <i> -> <ul>

        // check if list element has children and the first child isnt the NavAddSubBtn
        if(ListTag.children.length > 3 && !(ListTag.children[0].className.includes("navaddsubbtn"))) {
            let answerBool = confirm("Möchtest du die gesamte Kategorie inklusive der Unterkategorien löschen?");
            if(answerBool)  {
                ListTag.remove();
            }    
        } else {
            ListTag.parentElement.remove(); // <ul> -> <div>
        }
        let navIndexList = this.getNavIndexList();
        this.ReIndexNavItems(navIndexList); // assign id's to items
        this.applyEditButtons(); // remove and re-assign edit buttons
    }
}

class editAsideWorker {

    AsideMain = document.getElementsByTagName("aside")[0];
    AsideItem = {};
    AsideItemList = [];

    constructor() {
        this.AsideMain.children[0].insertAdjacentElement('beforeend', this.createAsideEditButton());
        this.getPreviousAside();
    }

    getPreviousAside() {
        let asideItems = document.querySelectorAll("aside button:not(#logoutBtn)");
        this.AsideItemList = [];

        for(let item of asideItems) {
            if(item.id != "editAsideBtn") {
                this.AsideItem = {};
                this.AsideItem.id = item.id.slice(1);
                this.AsideItem.title = item.innerHTML;
                this.AsideItemList.push(this.AsideItem);
            }
        }
    }

    createAsideEditButton() {
        var EditButton = document.createElement("button");
        EditButton.onclick = this.editAside.bind(this);
        EditButton.setAttribute("class","btn btn-danger float-left m-2")
        EditButton.setAttribute("id", "editAsideBtn");
        EditButton.textContent = "Edit"
        return EditButton
    }

    editAside() {
        document.getElementById("editNavBtn").disabled = true;
        document.getElementById("editFooterBtn").disabled = true;
        document.getElementById("editContentBtn").disabled = true;
        document.getElementById("editAsideBtn").remove();
        document.querySelectorAll("aside button").forEach((item) => item.setAttribute("data-em", "true"));
        this.AsideMain.children[0].insertAdjacentElement('beforeend', this.createAsideSaveButton());
        this.applyEditButtons();
    }

    createAsideSaveButton() {
        var SaveButton = document.createElement("button");
        SaveButton.onclick = this.saveAside.bind(this);
        SaveButton.setAttribute("class","btn btn-success float-right mr-4 mt-2")
        SaveButton.setAttribute("id", "saveAsideBtn");
        SaveButton.textContent = "Save"
        return SaveButton
    }

    saveAside() {
        document.getElementById("editNavBtn").disabled = false;
        document.getElementById("editFooterBtn").disabled = false;
        document.getElementById("editContentBtn").disabled = false;
        document.getElementById("saveAsideBtn").remove();
        document.querySelectorAll("aside button").forEach((item) => item.setAttribute("data-em", "false"));
        document.querySelectorAll("aside button").forEach(function(item) {
            if(item.id != "saveAsideBtn") {
                item.ondblclick = null;
				item.setAttribute("data-em", "false");
            }
        });
        document.querySelectorAll("aside i").forEach((item) => item.remove());
        this.AsideMain.children[0].insertAdjacentElement('beforeend', this.createAsideEditButton());
        this.getNextAside();
        this.createAsidePayload();
    }

    createAsidePayload() {
        var jsonObject = {"asideItems":this.AsideItemList};
        console.log(JSON.stringify(jsonObject));
        updateComponent("/5fzt78g4A7fnb882/aside", JSON.stringify(jsonObject), true);
        this.AsideItemList = [];
    }

    createAddAsideItemButton() {
        var AddButton = document.createElement("i");
        AddButton.setAttribute("class","bi bi-clipboard-plus asideaddbtn")
        AddButton.addEventListener('click', this.addAsideItem.bind(this));
        return AddButton
    }

    getAsideIndexList() {
        let asideItems = document.querySelectorAll("aside button:not(#logoutBtn)");
        let asideIndexList = []; // need an array here for splice (not applicable to navItems htmlcollection)
        asideItems.forEach(function(item, index) {
            if(item.id != "saveAsideBtn") {
                asideIndexList.push(index)
            }
        }); // get index  from current nav item
        return asideIndexList;
    }

    addAsideItem(e) {
        let asideIndexList = this.getAsideIndexList();
        let newAsideHtml = '<button data-em="true" class="aside-link btn w-75 m-1">Neues Item</button>';
        asideIndexList.splice(0, 0, 0); // insert new index
        e.target.insertAdjacentHTML('afterend', newAsideHtml); // insert new html
        this.ReIndexAsideItems(asideIndexList); // assign id's to items
        this.applyEditButtons(); // remove and re-assign edit buttons
    }

    ReIndexAsideItems(asideIndexList) {
        let asideItems = document.querySelectorAll("aside button");
        for(let i = 0; i < asideIndexList.length; i++) {
            if(asideItems[i].id != "saveAsideBtn") {
                asideItems[i].id = "a" + (i + 1); // ids beginning with 1
            }
        }
    }

    createDelAsideItemButton() {
        var DelButton = document.createElement("i");
        DelButton.setAttribute("class","bi bi-trash asidedelbtn")
        DelButton.addEventListener('click', this.delAsideItem.bind(this));
        return DelButton
    }
    
    delAsideItem(e) {
        let BtnTag = e.target.previousSibling.previousSibling;
        BtnTag.remove();
        let asideIndexList = this.getAsideIndexList();
        this.ReIndexAsideItems(asideIndexList); // assign id's to items
        this.applyEditButtons(); // remove and re-assign edit buttons
    }

    applyEditButtons() {
        var AsideItems = document.querySelectorAll("aside button:not(#logoutBtn)");
        var editBtns = document.querySelectorAll('aside i');
        editBtns.forEach((item) => item.remove()); // remove if any
        
        for(let item of AsideItems) {
            if(item.id != "saveAsideBtn") {
                if(!(item.id.slice(1) === '1')) {
                    item.insertAdjacentElement('afterend', this.createDelAsideItemButton());
                }
                item.insertAdjacentElement('afterend', this.createAddAsideItemButton());
                item.ondblclick = this.editAsideItem.bind(this);
            }
        }
    }

    editAsideItem(e) {
        e.stopPropagation();
        let formerValue = e.target.innerHTML;
        let editInput = document.createElement('input');
        editInput.type = "text";
        editInput.id = "editInput";
        editInput.size = "8";
        editInput.maxLength = "25";
        editInput.value = formerValue;
        editInput.onfocusout = this.saveAsideInput.bind(this);
        e.target.innerHTML = "";
        e.target.appendChild(editInput);
        document.getElementById('editInput').focus();
    }

    saveAsideInput(e) {
        let inputTag = document.getElementById('editInput');
        if(this.AsideMain.contains(inputTag)) {
            let inputValue = inputTag.value;
            inputTag.parentNode.innerHTML = inputValue;
        }
    }

    getAsideItemIndex() {
        for(let j=0; j < this.AsideItemList.length; j++) {
            if(this.AsideItemList[j].ID === this.AsideItem.id) {
                return j;
            }
        }
        return -1;
    }

    getNextAside() {
        let AsideList = document.querySelectorAll("aside button:not(#logoutBtn)");
        let i;
        this.AsideItemList = [];
        //AsideList.forEach((item, index) => console.log(index, item.id.slice(1), item.innerHTML));

        for(i=0; i < AsideList.length; i++) {
            if(AsideList[i].id != "editAsideBtn") {
                this.AsideItem = {};
                this.AsideItem.id = AsideList[i].id.slice(1);
                this.AsideItem.title = AsideList[i].innerHTML;
    
                var index = this.getAsideItemIndex();
                
                if(index === -1) {
                    this.AsideItemList.push(this.AsideItem);
                } else {
                    this.AsideItemList[index] = this.AsideItem;
                }
            }
        }        
        // this.AsideItemList.forEach((item, index) => console.log(index, item.ID, item.Title));
    }

}

class editFooterWorker {

    FooterMain = document.getElementsByTagName("footer")[0].children[0].children[0].children[0];
    FooterItem = {};
    FooterItemList = [];
    FooterOnClickList = [];

    constructor() {
        this.FooterMain.insertAdjacentElement('afterbegin', this.createFooterEditButton());
        this.getPreviousFooter();
    }

    getPreviousFooter() {
        let footerItems = document.querySelectorAll("footer a");

        for(let item of footerItems) {
                this.FooterItem = {};
                this.FooterItem.id = item.id.slice(1);
                this.FooterItem.title = item.innerHTML;
                this.FooterItemList.push(this.FooterItem);
        }
    }

    applyEditButtons() {
        var FooterItems = document.querySelectorAll("footer a");
        var editBtns = document.querySelectorAll('footer i');
        editBtns.forEach((item) => item.remove()); // remove if any
        for(let item of FooterItems) {
                if(!(item.id.slice(1) === '1')) {
                    item.insertAdjacentElement('afterend', this.createDelFooterItemButton());
                }
                item.insertAdjacentElement('afterend', this.createAddFooterItemButton());
                item.ondblclick = this.editFooterItem.bind(this);
        }
    }

    createFooterEditButton() {
        var EditButton = document.createElement("button");
        EditButton.onclick = this.editFooter.bind(this);
        EditButton.setAttribute("class","btn btn-danger float-left")
        EditButton.setAttribute("id", "editFooterBtn");
        EditButton.textContent = "Edit"
        return EditButton
    }

    editFooter() {
        document.getElementById("editNavBtn").disabled = true;
        document.getElementById("editAsideBtn").disabled = true;
        document.getElementById("editContentBtn").disabled = true;
        document.getElementById("editFooterBtn").remove();
        document.querySelectorAll("footer a").forEach((item) => item.setAttribute("data-em", "true"));
        this.FooterMain.insertAdjacentElement('beforeend', this.createFooterSaveButton());
        this.applyEditButtons();
    }

    editFooterItem(e) {
        e.stopPropagation();
        let formerValue = e.target.innerHTML;
        let editInput = document.createElement('input');
        editInput.type = "text";
        editInput.id = "editInput";
        editInput.size = "8";
        editInput.maxLength = "25";
        editInput.value = formerValue;
        editInput.onfocusout = this.saveFooterInput.bind(this);
        e.target.innerHTML = "";
        e.target.appendChild(editInput);
        document.getElementById('editInput').focus();
    }

    saveFooterInput(e) {
        let inputTag = document.getElementById('editInput');
        if(this.FooterMain.contains(inputTag)) {
            let inputValue = inputTag.value;
            inputTag.parentNode.innerHTML = inputValue;
        }
    }

    createFooterSaveButton() {
        var SaveButton = document.createElement("button");
        SaveButton.onclick = this.saveFooter.bind(this);
        SaveButton.setAttribute("class","btn btn-success ml-4")
        SaveButton.setAttribute("id", "saveFooterBtn");
        SaveButton.textContent = "Save"
        return SaveButton
    }

    saveFooter() {
        document.getElementById("editNavBtn").disabled = false;
        document.getElementById("editAsideBtn").disabled = false;
        document.getElementById("editContentBtn").disabled = false;
        document.getElementById("saveFooterBtn").remove();
        document.querySelectorAll("footer i").forEach((item) => item.remove());
        document.querySelectorAll("footer a").forEach((item) => {
			item.ondblclick = null;
			item.setAttribute("data-em", "false");
		});
        this.FooterMain.insertAdjacentElement('afterbegin', this.createFooterEditButton());
        this.getNextFooter();
        this.createFooterPayload();
    }

    createAddFooterItemButton() {
        var AddButton = document.createElement("i");
        AddButton.setAttribute("class","bi bi-plus-square footeraddbtn")
        AddButton.addEventListener('click', this.addFooterItem.bind(this));
        return AddButton
    }

    addFooterItem(e) {
        let footerIndexList = this.getFooterIndexList();
        let newFooterHtml = '<a data-em="false" role="button">Neuer Link</a>';
        footerIndexList.splice(0, 0, 0); // insert new index
        e.target.insertAdjacentHTML('afterend', newFooterHtml); // insert new html
        this.ReIndexFooterItems(footerIndexList); // assign id's to items
        this.applyEditButtons(); // remove and re-assign edit buttons
    }

    getFooterIndexList() {
        let footerItems = document.querySelectorAll("footer a");
        let FooterItemList = []; // need an array here for splice (not applicable to navItems htmlcollection)
        footerItems.forEach((_, index) => FooterItemList.push(index)); // get index  from current nav item
        return FooterItemList;
    }

    ReIndexFooterItems(FooterIndexList) {
        let footerItems = document.querySelectorAll("footer a");
        footerItems.forEach((item, index) => console.log("debug", item,index));
        for(let i = 0; i < FooterIndexList.length; i++) {
            footerItems[i].id = "f" + (i + 1); // ids beginning with 1
        }
    }

    createDelFooterItemButton() {
        var DelButton = document.createElement("i");
        DelButton.setAttribute("class","bi bi-dash-square footerdelbtn")
        DelButton.addEventListener('click', this.delFooterItem.bind(this));
        return DelButton
    }

    delFooterItem(e) {
        let FooterTag = e.target.previousSibling.previousSibling;
        FooterTag.remove();
        let footerIndexList = this.getFooterIndexList();
        this.ReIndexFooterItems(footerIndexList); // assign id's to items
        this.applyEditButtons(); // remove and re-assign edit buttons
    }

    getNextFooter() {
        let FooterList = document.querySelectorAll("footer a");
        let i;
        this.FooterItemList = [];
        //FooterList.forEach((item, index) => console.log(index, item.id.slice(1), item.innerHTML));

        for(i=0; i < FooterList.length; i++) {
                this.FooterItem = {};
                this.FooterItem.id = FooterList[i].id.slice(1);
                this.FooterItem.title = FooterList[i].innerHTML;
    
                var index = this.getFooterItemIndex();
                
                if(index === -1) {
                    this.FooterItemList.push(this.FooterItem);
                } else {
                    this.FooterItemList[index] = this.FooterItem;
                }
        }        
        // this.FooterItemList.forEach((item, index) => console.log(index, item.ID, item.Title));
    }

    getFooterItemIndex() {
        for(let j=0; j < this.FooterItemList.length; j++) {
            if(this.FooterItemList[j].ID === this.FooterItem.id) {
                return j;
            }
        }
        return -1;
    }

    createFooterPayload() {
        var jsonObject = {"footerItems":this.FooterItemList};
        console.log(JSON.stringify(jsonObject));
        updateComponent("/5fzt78g4A7fnb882/footer", JSON.stringify(jsonObject), false);
        this.FooterItemList = [];
    }

}


// credits to https://gist.github.com/Daniel-Hug/de9a165a6d9c74686854
// to unwrap necessary elements
function unwrap(wrapper) {
	// place childNodes in document fragment
	var docFrag = document.createDocumentFragment();
	while (wrapper.firstChild) {
		var child = wrapper.removeChild(wrapper.firstChild);
		docFrag.appendChild(child);
	}

	// replace wrapper with document fragment
	wrapper.parentNode.replaceChild(docFrag, wrapper);
}

// updates a specific component via fetch post
// target location is determined by path and data a.k.a. payload is already stringified json
function updateComponent(path, data, reload) {
    fetch(path, {
        method: "POST",
        headers: {
            'Accept': 'application/json, text/plain',
            'Content-Type': 'application/json;charset=UTF-8'
        },
        body: data
    })
    .then(response => {
        if (!response.ok) {
            return response.text().then(text => {throw new Error(text)})
        }
      })
      .catch(err => {
        console.log(err);
      });
    console.log("Settings have been saved.");
    if(reload) {
        setTimeout(function() {
            location.reload();
        }, 1000);
    }
}

function ContentonClickAddLabel() {
    var checkbox = document.getElementById("extSwitch");
    if (checkbox.checked == true) {
        var placeHolder = document.getElementsByClassName("jodit-placeholder");
        if(placeHolder[0]){
            placeHolder[0].remove();
        }
        document.getElementsByClassName("jodit-wysiwyg")[0].innerHTML = "<p>Sie werden weitergeleitet zu ...&nbsp;</p>"
    } else {
        document.getElementsByClassName("jodit-wysiwyg")[0].innerHTML = "<p></p>"
    }
}

function SignalLogout() {
    fetch("/5fzt78g4A7fnb882/logout", {
        method: 'POST'
    }).then(response => {
        if (!response.ok) {
            return response.text().then(text => {throw new Error(text)})
        }
      })
      .catch(err => {
        console.log(err);
      });
      document.getElementById("logoutBtn").remove();
      setTimeout(function() {
          location.reload();
      }, 500);
}

function parseStyleAttribute(e) {
    var x = e.getAttribute("style");
    var split = x.split(";");
    var xsplit = {};
    var xsplitS = {Value:[]}; // pointer replacement https://stackoverflow.com/questions/10231868/pointers-in-javascript
    var z;

    for (var j = 0; j < split.length; j++) {
        if (split[j] != "") {
            xsplit = {};
            z = split[j].split(": ");
            xsplit.key = z[0];
            xsplit.value = z[1];
            //console.log("xsplit", xsplit);
            xsplitS.Value.push(xsplit);
        }
    }

    for (var y = 0; y < xsplitS.Value.length; y++) {
        //console.log("xkey", xsplitS.Value[y].key);
        //console.log("xvalue", xsplitS.Value[y].value);
        switch(xsplitS.Value[y].key) {
            case "text-align":
                switch(xsplitS.Value[y].value) {
                    case "center":
                        if (!e.classList.contains("text-center")) {
                            e.classList.add("text-center");
                        }
                        if (e.classList.contains("text-right")) {
                            e.classList.remove("text-right");
                        }
                        if (e.classList.contains("text-left")) {
                            e.classList.remove("text-left");
                        }
                        break;
                    case "left":
                        if (!e.classList.contains("text-left")) {
                            e.classList.add("text-left");
                        }
                        if (e.classList.contains("text-center")) {
                            e.classList.remove("text-center");
                        }
                        if (e.classList.contains("text-right")) {
                            e.classList.remove("text-right");
                        }
                        break;
                    case "right":
                        if (!e.classList.contains("text-right")) {
                            e.classList.add("text-right");
                        }
                        if (e.classList.contains("text-center")) {
                            e.classList.remove("text-center");
                        }
                        if (e.classList.contains("text-left")) {
                            e.classList.remove("text-left");
                        }
                        break;
                    default:
                        break;
                }
            case "width":
                e.setAttribute("width", xsplitS.Value[y].value.replace("px",""));
                break;
            case "height":
                e.setAttribute("height", xsplitS.Value[y].value.replace("px",""));
                break;
            default:
                break;
        }
    }
}