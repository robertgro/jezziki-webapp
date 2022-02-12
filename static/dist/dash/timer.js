// credits https://stackoverflow.com/questions/20618355/the-simplest-possible-javascript-countdown-timer
/* is the old one, maybe merge this with the local storage timer later
function startTimer(duration, display) {
    var timer = duration, minutes, seconds;
    setInterval(function () {
        minutes = parseInt(timer / 60, 10);
        seconds = parseInt(timer % 60, 10);

        minutes = minutes < 10 ? "0" + minutes : minutes;
        seconds = seconds < 10 ? "0" + seconds : seconds;

        display.textContent = minutes + ":" + seconds;

        if (--timer < 0) {
            timer = duration;
        }
    }, 1000);
}

function initTimer() {
    var oneHour = 60 * 60,
    display = document.querySelector('#time');
    startTimer(oneHour, display);
}
*/

function cTimer() {
    var time, min, sec, minStr, secStr;

    if(typeof localStorage.getItem("min") !== 'undefined' && typeof localStorage.getItem("sec") !== 'undefined'){
        if(localStorage.getItem("min")!= null && localStorage.getItem("sec")!= null) {
            min = localStorage.getItem("min");
            sec = localStorage.getItem("sec");
            minStr = min.toString();
            secStr = "0" + sec.toString();
        } else {
            window.localStorage.clear();
            min = 10;
            sec = 0;
            minStr = min.toString();
            secStr = "0" + sec.toString();
            localStorage.setItem("min", min);
            localStorage.setItem("sec", sec);
        }
        //console.log(min, sec);
    }

    if(localStorage.getItem("min") == 0 && localStorage.getItem("sec") == 0) {
        window.localStorage.clear();
        min = 10;
        sec = 0;
        minStr = min.toString();
        secStr = "0" + sec.toString();
        localStorage.setItem("min", min);
        localStorage.setItem("sec", sec);
    }

    setInterval(function()
    {
        localStorage.setItem("min", min);
        localStorage.setItem("sec", sec);

        //console.log(min, sec, minStr, secStr);
        time=minStr+":"+secStr;
        document.getElementById("timer").innerHTML = time;
        
        if(sec == 00)
        {
            if(min !=0)
            {
                min--;
                sec=59;
                if(min < 10)
                {
                    minStr="0"+min.toString();
                }
                else
                {
                    minStr=min.toString();
                }
            }
        }
        else
        {
            sec--;
            if(sec < 10)
            {
                secStr="0"+sec.toString();
            }
            else
            {
                secStr=sec.toString();
            }
        }
    },1000);
}