customElements.define("notification-ele", class Notiication extends HTMLElement{
    constructor(){
        super()
        this.type = this.getAttribute("type")
        this.message = this.getAttribute("msg")
        this.warning = `
            <svg
                style="width: 80%;height: 80%;"
                width="64"
                height="63.999996"
                viewBox="0 0 16.933334 16.933332">
                <g
                    transform="translate(-87.307433,-97.949631)">
                <rect
                    style="fill:var(--warning_border_color);stroke-width:0.491643"
                    width="2.4377697"
                    height="10.683162"
                    x="94.408279"
                    y="98.478798"
                    ry="0.30981666" />
                <circle
                    style="fill:var(--warning_border_color);stroke-width:1.67205"
                    cx="95.627167"
                    cy="112.4334"
                    r="1.9196181" />
                </g>
            </svg>
        `
        this.success = `
           <svg
                style="width: 80%;height: 80%;"
                width="64.000008"
                height="57.549763"
                viewBox="0 0 16.933335 15.226707">
                <g
                    transform="translate(-46.262468,-78.805815)">
                <path
                    style="fill:transparent;stroke:var(--success_border_color);stroke-width:3.361;stroke-opacity:1"
                    d="m 47.390491,85.138996 c 1.894167,1.715301 3.479146,3.648796 4.407963,5.488994 3.713518,-5.874505 7.072233,-8.107992 10.440608,-10.440608"
                    id="checkmark"/>
                </g>
           </svg>
        `
        this.error = `
<svg
        style="width:70%;height:70%"
   width="63.999996"
   height="63.999996"
   viewBox="0 0 16.933332 16.933332">
  <g
     transform="translate(-51.903267,-103.62282)">
    <rect
       style="fill:var(--error_border_color);stroke-width:0.0928013"
       id="rect113"
       width="3.4994376"
       height="20.996624"
       x="120.19751"
       y="26.072937"
       ry="1.7497188"
       transform="rotate(45)" />
    <rect
       style="fill:var(--error_border_color);stroke-width:0.0928013"
       id="rect273"
       width="3.4994376"
       height="20.996624"
       x="-38.320965"
       y="111.44891"
       ry="1.7497188"
       transform="rotate(-45)" />
  </g>
</svg>
        `

        this.innerHTML = `
            <div id="alert" class="${this.type === "warning" ? "alert_warning" : this.type === "success" ? "alert_success" : "alert_error"}">
                <div class="alert_circle">
                ${this.type === "warning" ? this.warning : this.type === "success" ? this.success : this.error}
                </div>
                <p>${this.message}</p>
            </div>
        `
    }
})
