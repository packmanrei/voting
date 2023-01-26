var numOfChoice = 3;
function addChoice(){

    // create choice field
    var choice = document.createElement('input');
    choice.type = 'text';
    choice.name = 'choice';
    choice.className = 'post_choice';
    choice.placeholder = '選択肢' + numOfChoice;
    
    // create br
    var br =  document.createElement('br');

    // select place
    var choice_box = document.querySelector('.choice_box');
    
    // put elements
    choice_box.appendChild(choice);
    choice_box.appendChild(br);
    numOfChoice++;
}


function addChoiceID(){
    document
    .querySelectorAll("button[name='button']")
    .forEach((btn, i) => {
      btn.value = i + 1;
    });
}

function lastCheck(){
  var btn = querySelector('.delete_button');
  btn.addEcentListener('click', function(){
    var check = window.confirm('本当によろしいですか？');
    if(result){
      location.href = ''
    }
  })
  var check = window.confirm
}


function callConditions(){
	const conditions = document.querySelector('.conditions');

	if(conditions.style.display=="block"){
		conditions.style.display ="none";
	}else{
		conditions.style.display ="block";
  }
}