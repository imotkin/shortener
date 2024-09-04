function copyText() {
    const link = document.getElementsByTagName("a")[0];   
    navigator.clipboard.writeText(link.textContent);

    const notification = document.getElementById("copy-notification");
    notification.classList.add("show");
    
    setTimeout(() => {
        notification.classList.remove("show");
    }, 5000);
}