const formatDate = (isoString) => {
  if (!isoString) return "-";
  const date = new Date(isoString);
  return date.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  });
};