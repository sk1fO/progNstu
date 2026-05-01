import spinnerImage from '../assets/spinner.png';
import './Spinner.css';

function Spinner({ fullPage = true }) {
  const spinner = (
    <div className="spinner-container">
      <img src={spinnerImage} alt="Загрузка" className="spinner-image" />
    </div>
  );

  if (fullPage) {
    return (
      <div className="spinner-fullpage">
        {spinner}
        <p>Мадуро-СПИННЕР-грузит</p>
      </div>
    );
  }

  return spinner;
}

export default Spinner;