const text = {
    title: "逸雨网络",
    subtitle: "逸雨网络",
    copyright: "2026 逸雨网络. All rights reserved.",
    contact: "接程序接口定制 UI设计 漏洞修复 联系：QQ 204325164",
    quote: "< Stay hungry, stay foolish. – Steve Jobs >",
    antiCopy: "Stop! 抄袭固然是快 但是我想cnmb!",
    warning1: "这里是给开发者准备的调试工具。",
    warning2: "如果有人让你在这里复制粘贴任何代码，",
    warning3: "他 99% 是在尝试盗取你的账号 刷你的账户余额。",
    warning4: "除非你非常清楚自己在做什么，否则请立即关闭控制台。"
}

console.log(`%c${text.title}`, `
    font-size: 48px;
    font-weight: bold;
    color: #007AFF;
    text-shadow: 2px 2px 4px rgba(0,122,255,0.3);
`);

console.log(`%c${text.subtitle}`, `
    font-size: 24px;
    background: linear-gradient(90deg, #007AFF, #5AC8FA);
    color: white;
    padding: 6px 16px;
    border-radius: 20px;
    display: inline-block;
    margin: 10px 0;
`);

console.log(`%c${text.copyright}`, `
    font-size: 16px;
    color: #666;
    margin: 5px 0;
`);

console.log(`%c${text.contact}`, `
    font-size: 18px;
    color: #007AFF;
    font-weight: bold;
    margin: 8px 0;
`);

console.log(`%c${text.quote}`, `
    font-size: 16px;
    color: #999;
    font-style: italic;
    margin: 8px 0;
`);

console.log(`%c${text.antiCopy}`, `
    font-size: 22px;
    color: #ff3333;
    font-weight: bold;
    margin: 15px 0;
`);

console.log(`%c${text.warning1}`, `
    font-size: 16px;
    color: #333;
`);
console.log(`%c${text.warning2}`, `
    font-size: 16px;
    color: #333;
`);
console.log(`%c${text.warning3}`, `
    font-size: 16px;
    color: #333;
`);
console.log(`%c${text.warning4}`, `
    font-size: 16px;
    color: #333;
`);

console.log(`%c${"─".repeat(60)}`, `
    color: #007AFF;
    font-size: 20px;
`)