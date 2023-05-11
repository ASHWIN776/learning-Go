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
        title = "",
        html = "",
        callback = () => {},
        willOpen = () => {},
        didOpen = () => {},
        preConfirm = () => {}
    }){
        const { value: formValues } = await Swal.fire({
            title,
            html,
            focusConfirm: false,
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
    notie.alert({
    type: type,
    text: msg,
    })
}