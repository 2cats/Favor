var summer;
var progress;
var PAGE_SIZE = 10;
var pagecount;
var filter = ""
var ajaxList = []

$('#fileupload').on('change', function prepareUpload(event) {
    uploadFiles(event.target.files);
});
var FileButton = function (context) {
    var ui = $.summernote.ui;
    var button = ui.button({
        contents: 'Upload',
        tooltip: 'Upload Files',
        click: function () {
            $('#fileupload').trigger('click');
        }
    });
    return button.render(); // return button as jquery object 
}
var SubmitButton = function (context) {
    var ui = $.summernote.ui;
    var button = ui.button({
        contents: '<div style="color:green;">Submit</>',
        tooltip: 'Submit',
        click: function () {
            var jsonData = {
                op: "INSERT",
                data: summer.summernote('code')
            }
            ajax("msg", "POST", JSON.stringify(jsonData), updateMsg, false, $('#msg'))
        }
    });
    return button.render(); // return button as jquery object 
}
var summerSubmitConf = {
    toolbar: [
        // [groupName, [list of button]]
        ['style', ['bold', 'italic', 'underline', 'clear']],
        ['font', ['strikethrough', 'superscript', 'subscript']],
        ['fontsize', ['fontsize']],
        ['color', ['color']],
        ['para', ['ul', 'ol', 'paragraph']],
        ['height', ['height']],
        ['codeview', ['codeview']],
        ['file', ['file']],
        ['submit', ['submit']]
    ],
    codemirror: { // codemirror options
        theme: 'monokai'
    },
    minHeight: '300px',
    buttons: {
        file: FileButton,
        submit: SubmitButton
    }
};

$('ul.nav > li').click(function (e) {
    // e.preventDefault();
    if ($(this).hasClass("active")) {
        $('ul.nav > li').removeClass('active');
    } else {
        $('ul.nav > li').removeClass('active');
        $(this).addClass('active');
    }
    if ($("#editNav").hasClass("active")) {
        $('html,body').animate({
            scrollTop: '0px'
        }, 200)
        $('#outersummer').show(600)
    } else {
        $('#outersummer').hide(400)
    }
});
$('#filter').blur(function (e) {
    e.preventDefault();
    changeFilter();
})
$('#filter').keydown(function (e) {
    if (e.keyCode == 13) {
        e.preventDefault();
        changeFilter();
    }
});

function changeFilter() {
    filter = $('#filter').val()
    console.log(ajaxList)
    for (var k in ajaxList) {
        ajaxList[k].abort();
    }
    ajaxList = []
    updateMsg(1)
}

function currentPage() {
    return parseInt($('#currpage').text().replace(/[^0-9]/ig, ""));
}

function onNavLeftClick() {
    updateMsg(currentPage() - 1)
}

function onNavRightClick() {
    updateMsg(currentPage() + 1)
}

function onUploadSuccess(files) {
    for (var filename in files) {
        var ldot = files[filename].lastIndexOf(".");
        var type = files[filename].substring(ldot + 1);
        var node;
        if (type === "jpg" || type === "png" || type === "gif" || type === "bmp") {
            node = document.createElement('img');
            node.setAttribute("src", files[filename]);
            node.innerHTML = filename
        } else {
            node = document.createElement('a');
            node.setAttribute("href", files[filename]);
            node.innerHTML = filename
        }
        $('#summernote').summernote('insertNode', node);
    }
}

function ajax(url, type, data, onsuccess, onprogress, showui) {

    if (showui) {
        $(showui).html("<div id='loadingui' class='load loading-2'><span></span><span></span><span></span><span></span><span></span></div>")
    }
    var newajax = $.ajax({
        url: url,
        type: type,
        data: data,
        cache: false,
        dataType: 'json',
        processData: false, // Don't process the files
        contentType: false, // Set content type to false as jQuery will tell the server its a query string request
        success: function (data, textStatus, jqXHR) {
            if (typeof data.error === 'undefined' && data.status === 'SUCCESS') {
                dhtmlx.message({
                    text: "Operation Success",
                    type: "success"
                });
                if (onsuccess) {
                    onsuccess(data.data)
                }
                console.log("SUCCESS")
            } else {
                dhtmlx.message({
                    text: "Request " + url + " Error\n",
                    type: "error"
                });
             }

        },
        error: function (jqXHR, textStatus, errorThrown) {
            dhtmlx.message({
                text: "Request " + url + " Error\n",
                type: "error"
            });
        },
        progress: function (e) {
            //make sure we can compute the length
            if (e.lengthComputable) {
                if (onprogress) {
                    var pct = (e.loaded / e.total) * 100;
                    onprogress(pct)
                }
            } else {
                console.warn('Content Length not reported!');
            }
        },
        complete: function (e) {
            $(showui).children('#loadingui').remove()
        }
    });
    ajaxList.push(newajax)

}



function uploadFiles(files) {
    var data = new FormData();
    $.each(files, function (key, value) {
        data.append(key, value);
    });
    ajax("postfiles", "POST", data,
        function (data) {
            onUploadSuccess(data)
        },
        function (pct) {
            progress.set(pct)
            if (pct == 100) {
                progress.set(0);
            }
        }
    );
}

function fetchMsg(nth) {
    if (!nth) {
        nth = 1;
    }
    $("#msg").children("*").remove();
    var url = "msg?pagesize=" + PAGE_SIZE + "&nth=" + nth + "&filter=" + filter;
    ajax(url, "GET", "", function (data) {
        $("#msg").children("*").remove()
        for (var i in data) {
            var openDiv = $("<div class='msgitem textarea'></div>")
            $(openDiv).attr('id', data[i].Id);
            // var delbtn=$("<div class='btngroup'><button type='button' class='editbutton btn  btn-xs btn-success'>&nbsp;&nbsp;Edit&nbsp;&nbsp;</button> <button type='button' class='btn  btn-xs btn-danger'>Delete</button></div>")
            var delbtn = $("<div class='btngroup'></div>")
            var idLabel = $("<span class='idlabel'></span>")
            var dd = $("<button type='button' class='delbtn btn  btn-xs btn-danger'>Delete</button>")
            delbtn.hide();
            openDiv.addClass("msg")
            openDiv.html(data[i].Content);
            idLabel.html(data[i].Id)
            delbtn.appendTo(openDiv)
                // idLabel.appendTo(delbtn)
            dd.appendTo(delbtn)
            $(openDiv).hide().appendTo($("#msg")).show()
        }

        $(".msgitem").mouseenter(function () {
            console.log("over");
            $(this).children(".btngroup").fadeIn(400);
        });
        $(".msgitem").mouseleave(function () {
            $(this).children(".btngroup").hide();
        });
        $('.delbtn').click(function (e) {
            var msg = $(this).parent().parent();
            msg.removeClass("textarea");
            msg.hide(300)
            var jsonData = {
                op: "DELETE",
                data: msg.attr('id')
            }
            ajax("msg", "POST", JSON.stringify(jsonData))
        })
        $('#currpage').html("Page " + nth + " <span class='caret'></span>")

        if (nth >= pagecount) {
            $('#rightnav').addClass('disabled');
        } else {
            $('#rightnav').removeClass('disabled');
        }
        if (nth <= 1) {
            $('#leftnav').addClass('disabled');
        } else {
            $('#leftnav').removeClass('disabled');
        }
        var imgs = $(".msgitem img");
        var divWidth = parseFloat($("#msg").css("width"))
        for (var k in imgs) {
            var xradio = parseFloat($(imgs[k]).css("height")) / parseFloat($(imgs[k]).css("width"));
            $(imgs[k]).css("width", divWidth * 0.4 + "px");
            $(imgs[k]).css("height", divWidth * 0.4 * xradio + "px");
        }
    }, function (pct) {
        progress.set(pct)
        if (pct == 100) {
            progress.set(0);
        }
    }, $('#msg'))
}

function onNavClick(num) {
    console.log(num);
    updateMsg(num);
}

function updateMsg(nth) {
    var total = 0;
    url = "msg?filter=" + filter
    ajax(url, "GET", "", function (data) {
        total = data
        console.log(total)
        pagecount = Math.ceil(total / PAGE_SIZE)
        console.log(pagecount)
        $('#menunav').children("*").remove()
        for (var i = 1; i <= pagecount; i++) {
            var node = $("<li></li>")
            $(node).html("<a href='javascript:void(0)' onclick='onNavClick(" + i + ");'>" + i + "</a>")
            $(node).appendTo("#menunav")
        }
        fetchMsg(nth)

    }, false, $('#msg'))
}
$(document).ready(function () {
    console.log("ready")
    dhtmlx.message.position = "bottom"; // possible values "top" or "bottom"
    $('#outersummer').hide()
    summer = $('#summernote').summernote(summerSubmitConf)
    progress = progressJs('#progress').start()
    progress.set(0)
    updateMsg(1)
});
