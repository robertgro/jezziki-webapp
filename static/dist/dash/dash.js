var target = document.querySelector('#root');
var timerHTML = '<span class="text-danger font-weight-bold fixed-bottom ml-2 mb-2" id="timer">10:00</span>';
var logoutHTML = '<button class="btn btn-warning float-left m-3 fixed-top" id="logoutBtn">Logout</button>';
var contentWorkerObject, navWorkerObject, asideWorkerObject, footerWorkerObject;

function mutate(mutations) {
    mutations.forEach(function(mutation) {
      //console.log(mutation.type);
      //console.log(mutation.addedNodes);
      //console.log(mutation.removedNodes);

      if(mutation.type === "attributes") {
          if(mutation.attributeName === "style") {
            var n = mutation.target; // return node
            //console.log(n, n.getAttribute("style"), n.nodeName);
            if(n.nodeName === "P" || n.nodeName === 'IMG' || n.nodeName === 'IFRAME') {
                if(n.getAttribute("style") != null) {
                    parseStyleAttribute(n);
                }
            }
          }
      }
      
      // restores the content worker in case of new received post data
        mutation.removedNodes.forEach(function(removedNode) {
                //console.log(removedNode.id);
                if(removedNode.id === 'editContentBtn') {
                    var saveBtn = document.getElementById("saveContentBtn");
                    if(!saveBtn) {
                        contentWorkerObject = new editContentWorker();
                    }
                } else if(removedNode.id === 'editor') {
                    var saveBtn = document.getElementById("saveContentBtn");
                    if(!saveBtn) {
                        contentWorkerObject = new editContentWorker();
                    }
                }

                removedNode.childNodes.forEach(function(anode){
                    if(anode.id === 'logoutBtn') {
                    var logoutBtn = document.getElementById("logoutBtn");
                    if(!logoutBtn) {
                        let v = document.getElementsByTagName("aside");
                        v[0].querySelector('div').insertAdjacentHTML('beforeend', logoutHTML);
                        document.getElementById("logoutBtn").onclick = SignalLogout;
                    }
                }
            })
        })
    
        mutation.addedNodes.forEach(function(anode){
            //console.log(anode.nodeName);
            anode.childNodes.forEach(function(bnode) {
                //console.log(bnode.nodeName);
                bnode.childNodes.forEach(function(cnode) {
                    //console.log(cnode.nodeName);
                    cnode.childNodes.forEach(function(lnode) {
                        //console.log(lnode.nodeName);

                        // If we at this node, nodes have been added into dom successfully
                        if(lnode.nodeName === 'ASIDE') {
                            
                        contentWorkerObject = new editContentWorker();
                        navWorkerObject = new editNavWorker();
                        asideWorkerObject = new editAsideWorker();
                        footerWorkerObject = new editFooterWorker();

                        // Timer preparations
                        let h = document.getElementsByTagName('header');
                        h[0].querySelector('div > div > img').insertAdjacentHTML('afterend', timerHTML);

                        var logoutBtn = document.getElementById("logoutBtn");
                        if(!logoutBtn) {
                            let v = document.getElementsByTagName("aside");
                            v[0].querySelector('div').insertAdjacentHTML('beforeend', logoutHTML);
                            document.getElementById("logoutBtn").onclick = SignalLogout;
                        }

                        // unwrap a link
                        unwrap(document.querySelector('#root > a'));
                        // Insert Timer
                        cTimer();
                        }
                    });
                });
            });
        });
    });
  }

var observer = new MutationObserver(mutate);

var config = { attributes: true, childList: true, characterData: true, subtree: true };

observer.observe(target, config);