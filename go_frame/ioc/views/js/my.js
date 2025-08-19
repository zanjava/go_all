function show_login() {//展示要登录的样式
    $('#user_name').html("");
    document.querySelector('#log_out').style.display = 'none';
    document.querySelector('#login').style.display = 'block';
    let issueBtn = document.querySelector('#issue_bnt');
    if(issueBtn){
        issueBtn.style.display = 'none';
    }
    let deleteBtn = document.querySelector('#delete_bnt');
    if(deleteBtn){
        deleteBtn.style.display = 'none';
    }
    let editBtn = document.querySelector('#edit_bnt');
    if(editBtn){
        editBtn.style.display = 'none';
    }
};

function show_out(user_name) {//展示要退出的样式
    $('#user_name').html(user_name);
    document.querySelector('#log_out').style.display = 'block';
    document.querySelector('#login').style.display = 'none';
    let issueBtn = document.querySelector('#issue_bnt');
    if(issueBtn){
        issueBtn.style.display = 'block';
    }
    // let deleteBtn = document.querySelector('#delete_bnt');
    // if(deleteBtn){
    //     deleteBtn.style.display = 'block';
    // }
    // let editBtn = document.querySelector('#edit_bnt');
    // if(editBtn){
    //     editBtn.style.display = 'block';
    // }
};