const attention = Prompt()

function Prompt(){
    function toast({
    msg = "",
    icon = "success",
    position = "top-end"
    }){
        const Toast = Swal.mixin({
            toast: true,
            icon,
            position,
            title: msg,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
            toast.addEventListener('mouseenter', Swal.stopTimer)
            toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    function success({
    title = "",
    msg = "",
    }){
        Swal.fire({
            icon: 'success',
            title,
            text: msg,
        })
        }

    function error({
    title = "",
    msg = "",
    }){
        Swal.fire({
            icon: 'error',
            title,
            text: msg,
        })
    }

    async function custom({
        title = "",
        html = "",
        callback = () => {}
    }){
        const { value: formValues } = await Swal.fire({
            title,
            html,
            focusConfirm: false,
            showCancelButton: true,
            willOpen: () => {
                const elem = document.getElementById('reservation-dates-modal');
                const rangepicker = new DateRangePicker(elem, {
                    showOnFocus: true,
                    orientation: "auto top"
                }); 
            },
            preConfirm: () => {
                return [
                    document.getElementById('startDate').value,
                    document.getElementById('endDate').value
                ]
            },
            didOpen: () => {
            document.getElementById('startDate').removeAttribute("disabled"),
            document.getElementById('endDate').removeAttribute("disabled")
            }
        })

        if (formValues) {
            callback()
        }
    }

    return {
        toast,
        error,
        success,
        custom
    }
}