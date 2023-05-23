const attention = Prompt()

// Sweetalert prompts
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
        icon = "",
        title = "",
        html = "",
        showConfirmButton = true,
        callback = () => {},
        willOpen = () => {},
        didOpen = () => {},
        preConfirm = () => {}
    }){
        const { value: formValues } = await Swal.fire({
            icon,
            title,
            html,
            focusConfirm: false,
            showConfirmButton: showConfirmButton,
            showCancelButton: true,
            willOpen,
            preConfirm,
            didOpen
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

// Notie alert
function notieAlert(type, msg)
{
    console.log("Calling notie")
    notie.alert({
    type: type,
    text: msg,
    })
}