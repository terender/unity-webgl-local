let local = false
// 保存原始的fetch函数
const originalFetch = fetch;

// 重写fetch函数
window.fetch = function(input, init) {
  if (!local) {
    return originalFetch(input, init)
  }
  // 获取的资源地址不是字符串
  if (typeof input !== 'string') {
    return originalFetch(input, init)
  }
  // 如果本来就是 js 文件则使用原始加载方法
  if (input.endsWith('.js')) {
    return originalFetch(input, init)
  }
  // 获取的资源不是Unity资源文件
  let basePath = removeFileNameFromPath(window.location.href);
  let assetPath = decodeURI(decodeURI(input)).replace(basePath, '');
  if (!assetPath.startsWith('Build/') && 
      !assetPath.startsWith('StreamingAssets')) {
    return originalFetch(input, init)
  }

  // 你的自定义逻辑
  console.log(`读取Unity资源文件: ${assetPath}`);

  var varName = `_uinty_asset_` + assetPath.split('/').pop().replaceAll('.', '_').replaceAll('-', '_')
  var contentType = 'application/octet-stream'
  if (assetPath.endsWith('.wasm'))
    contentType = 'application/wasm'

  if (window[varName] !== undefined) {
    let asset = window[varName];
    const blob = new Blob([asset], {type: contentType});
    const response = new Response(blob, {
      status: 200, 
      headers: { 'Content-Type': contentType }
    });
    return new Promise(function(resolve, reject) { resolve(response) })
  }

  return new Promise(function(resolve, reject) {
    var script = document.createElement('script');
    script.src = assetPath + '.js';    
    script.onload = function() {
      let asset = window[varName];
      if (asset !== undefined) {
        const blob = new Blob([asset], {type: contentType});
        const response = new Response(blob, {
          status: 200, 
          headers: { 'Content-Type': contentType } 
        });
        // 返回Response对象
        resolve(response); 
      } else {
        // 返回错误
        reject(new Error(`${varName} is not defined`));
      }
    };

    script.onerror = function() {
      reject(new Error(`${assetPath} load error`));
    };

    document.body.appendChild(script);
  });
  
};

function removeFileNameFromPath(url) {
  // 使用URL对象解析当前URL
  const parsedUrl = new URL(url);

  // 获取路径名并分割成数组
  let pathSegments = parsedUrl.pathname.split('/');

  // 移除数组的最后一个元素（即文件名）
  // 检查最后一个元素是否为空，空通常意味着路径以斜杠结束
  if (pathSegments[pathSegments.length - 1]) {
    pathSegments.pop();
  }

  // 重新构建没有文件名的路径
  const newPathname = pathSegments.join('/').replace(/^\/+|\/+$/g, '');

  // 构建新的URL
  return decodeURI(`${parsedUrl.origin}/${newPathname}/`);
}