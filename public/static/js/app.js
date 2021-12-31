Vue.component('post', {
    props: ['item'],
    template: `
    <div class="post">
        <div class="meta">{{ item.user_preferred_name }} ({{ item.user_fingerprint }}) @ {{ item.timestamp }}</div>
        <div class="message">{{ item.message }}</div>
    </div>
    `
})

var app = new Vue({
    delimiters: ['$[', ']'],
    el: '#app',
    data: {
        posts: [
            { "id": "2154ce9b-38da-43b9-9f3f-ab4687cd5770", "timestamp": "2021-12-31T14:21:29.674408Z", "user_id": "cace32e0-7697-4981-9add-97ca2d482712", "user_preferred_name": "BarveyHirdmann", "user_fingerprint": "SHA256:U7zgVIf+5ypHTGTuw16LOgnUJ497K9rX+bpwnvXu6lg", "message": "This is a test message" }
        ]
    }
})