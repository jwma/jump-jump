import QRCode from 'qrcode'

export async function generateQRCodeDataURL(text: string, size = 256): Promise<string> {
  return QRCode.toDataURL(text, {
    width: size,
    margin: 2,
    color: { dark: '#000000', light: '#ffffff' },
  })
}

export function downloadQRCode(dataURL: string, filename: string) {
  const link = document.createElement('a')
  link.download = filename
  link.href = dataURL
  link.click()
}
