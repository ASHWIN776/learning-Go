(() => {
    'use strict'

    // Fetch all the forms we want to apply custom Bootstrap validation styles to
    const forms = document.querySelectorAll('.needs-validation')

    // Loop over them and prevent submission
    Array.from(forms).forEach(form => {
    form.addEventListener('submit', event => {
        if (!form.checkValidity()) {
        event.preventDefault()
        event.stopPropagation()
        }

        form.classList.add('was-validated')
    }, false)
    })
})()

const elem = document.getElementById('reservation-dates');
const rangepicker = new DateRangePicker(elem, {
    format: "dd-mm-yyyy"
}); 

function alert(type, msg)
{
    notie.alert({
    type: type,
    text: msg,
    })
}

const alertButton = document.querySelector("#alertButton")
const toastButton = document.querySelector("#toastButton")
const successModalButton = document.querySelector("#successModalButton")
const errorModalButton = document.querySelector("#errorModalButton")
const formModalButton = document.querySelector("#formModalButton")

alertButton.addEventListener("click", () => {
    alert("success", "Thie is the alert message")
})

toastButton.addEventListener("click", () => {
    attention.toast({msg:"Hello World"})
})

successModalButton.addEventListener("click", () => {
    attention.success({msg: "This is a success message", title:"Authenticated"})
})

errorModalButton.addEventListener("click", () => {
    attention.error({msg: "This is an error message", title: "Auth Error"})
})

formModalButton.addEventListener("click", () => {
    const html = `
    <form class="needs-validation" id="reservation-form" novalidate>
        <div class="row w-90 mx-auto my-4" id="reservation-dates-modal">
        <div class="col">
            <input disabled id="startDate" placeholder="Arrival" class="form-control" required type="text" name="startDate">
        </div>
        <div class="col">
            <input disabled id="endDate" placeholder="Departure" class="form-control" required type="text" name="endDate">  
        </div>
        </div>
    </form> 
    `
    attention.custom({title: "Book Rooms", html})
})